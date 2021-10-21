package go_go_github_badge

import (
	"fmt"
	"net/http"
	"sort"
	"strings"

	"github.com/gin-gonic/gin"
)

const cacheBadgeSec = 86400 * 7
const contributionDays = 30

var allowedUsers = map[string]bool{}

func SetAllowedUsers(allowed []string) {
	for _, u := range allowed {
		allowedUsers[strings.TrimSpace(u)] = true
	}
}

func Run() {
	CreateClient()

	r := gin.Default()

	r.Static("/css", "./static/css")
	r.Static("/image", "./static/image")
	r.StaticFile("/crossdomain.xml", "./static/crossdomain.xml")
	r.LoadHTMLGlob("templates/*.gohtml")

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	r.GET("/badge/:username", generateBadge)

	// port 8080
	r.Run()
}

func isUserAllowed(username string) bool {
	return allowedUsers[username]
}

func generateBadge(c *gin.Context) {
	username := c.Param("username")

	if !isUserAllowed(username) {
		c.String(http.StatusForbidden, "username not allowed")
		return
	}

	user, err := GetUser(username)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	contributions, err := GetLatestContributions(user, contributionDays)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	// latest project
	byRepo := contributions.User.ContributionsCollection.CommitContributionsByRepository
	sort.Slice(byRepo, func(i, j int) bool {
		return byRepo[j].Repository.UpdatedAt.Before(byRepo[i].Repository.UpdatedAt.Time)
	})

	forkCount, err := GetForkCount(username)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	repoStats, err := GetRepoStats(username)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	// fmt.Println("most recent: ", github.Stringify(byRepo[0].Repository.UpdatedAt.String()))
	// fmt.Println(github.Stringify(byRepo))
	// fmt.Println(github.Stringify(contributions))
	// fmt.Println(github.Stringify(contributions.User.Repositories.Nodes[0]))

	langs := []string{}
	for _, r := range repoStats.Languages {
		langs = append(langs, r.Name)
	}

	c.Header("cache-control", fmt.Sprintf("public, max-age=%d", cacheBadgeSec))

	c.HTML(http.StatusOK, "badge.gohtml", gin.H{
		"username":           username,
		"title":              "Main website",
		"User":               user,
		"Followers":          contributions.User.Followers.TotalCount,
		"TotalContributions": contributions.User.ContributionsCollection.ContributionCalendar.TotalContributions,
		"Days":               contributionDays,
		"TotalRepos":         contributions.User.Repositories.TotalCount,
		"Repos":              contributions.User.Repositories.TotalCount - forkCount,
		"Forks":              forkCount,
		"Stargazers":         repoStats.StargazerCount,
		"AllLanguages":       strings.Join(langs, ", "),
		"TopLanguages":       strings.Join(langs[0:3], ", "),
	})
}
