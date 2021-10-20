package go_go_github_badge

import (
	"fmt"
	"net/http"
	"sort"

	"github.com/gin-gonic/gin"
	"github.com/google/go-github/v39/github"
)

func Run() {
	CreateClient()

	r := gin.Default()

	r.Static("/js", "./static/js")
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

func generateBadge(c *gin.Context) {
	username := c.Param("username")
	user, err := GetUser(username)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	contributions, err := GetLatestContributions(user)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	// latest project
	byRepo := contributions.User.ContributionsCollection.CommitContributionsByRepository
	sort.Slice(byRepo, func(i, j int) bool {
		fmt.Printf("%s: %s\n", byRepo[i].Repository.Name, byRepo[i].Repository.UpdatedAt.String())
		return byRepo[j].Repository.UpdatedAt.Before(byRepo[i].Repository.UpdatedAt.Time)
	})

	fmt.Println("most recent: ", github.Stringify(byRepo[0].Repository.UpdatedAt.String()))
	fmt.Println(github.Stringify(byRepo))
	fmt.Println(github.Stringify(contributions))

	fmt.Println(github.Stringify(contributions.User.Repositories.Nodes[0]))

	c.HTML(http.StatusOK, "badge.gohtml", gin.H{
		"username":           username,
		"title":              "Main website",
		"User":               user,
		"TotalContributions": contributions.User.ContributionsCollection.ContributionCalendar.TotalContributions,
		"Days":               30,
	})
}
