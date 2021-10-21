package go_go_github_badge

import (
	"context"
	"os"
	"sort"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/go-github/v39/github"
	"github.com/pkg/errors"
	"github.com/shurcooL/githubv4"
	"golang.org/x/oauth2"
)

var client = github.NewClient(nil)
var gclient *githubv4.Client

func CreateClient() {
	src := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: os.Getenv("GITHUB_TOKEN")},
	)
	httpClient := oauth2.NewClient(context.Background(), src)
	gclient = githubv4.NewClient(httpClient)
	// gclient = githubv4.NewEnterpriseClient("https://graphql.github.com/graphql/proxy", httpClient)
}

func GetUser(username string) (*github.User, error) {
	user, _, err := client.Users.Get(context.Background(), username)
	if err != nil {
		return nil, errors.Wrap(err, "error fetching user info")
	}
	return user, nil
}

type LatestContributions struct {
	User struct {
		Followers struct {
			TotalCount int
		}

		Repositories struct {
			TotalCount int
			Nodes      []struct {
				Name      string
				IsPrivate bool
				PushedAt  githubv4.DateTime

				DefaultBranchRef struct {
					Name   string
					Target struct {
						SpreadCommits struct {
							History struct {
								TotalCount int
							} `graphql:"history(since: $since)"`
						} `graphql:"... on Commit"`
					}
				}
			}
		} `graphql:"repositories(first: 10, orderBy: {field: PUSHED_AT, direction: DESC})"`

		ContributionsCollection struct {
			TotalCommitContributions int
			ContributionCalendar     struct {
				TotalContributions int
			}
			CommitContributionsByRepository []struct {
				Repository struct {
					UpdatedAt githubv4.DateTime
					Name      string
					URL       string
					IsPrivate bool
				}
				Contributions struct {
					TotalCount int
				}
			}
		} `graphql:"contributionsCollection(from: $from, to: $to)"`
	} `graphql:"user(login: $login)"`
}

func GetLatestContributions(user *github.User, days int) (*LatestContributions, error) {
	to := time.Now()
	from := to.Add(-time.Hour * 24 * time.Duration(days))
	args := gin.H{
		"login": githubv4.String(*user.Login),
		"from":  githubv4.DateTime{from},
		"to":    githubv4.DateTime{to},
		"since": githubv4.GitTimestamp{from},
	}

	query := LatestContributions{}
	err := gclient.Query(context.Background(), &query, args)
	if err != nil {
		return nil, errors.Wrap(err, "error fetching contributions")
	}

	return &query, nil
}

type ForkCount struct {
	User struct {
		Repositories struct {
			TotalCount int
		} `graphql:"repositories(isFork: true)"`
	} `graphql:"user(login: $login)"`
}

func GetForkCount(login string) (int, error) {
	args := gin.H{
		"login": githubv4.String(login),
	}
	query := ForkCount{}
	err := gclient.Query(context.Background(), &query, args)
	if err != nil {
		return 0, errors.Wrap(err, "error fetching fork count")
	}

	return query.User.Repositories.TotalCount, nil
}

type RepoInfo struct {
	StargazerCount  int
	IsFork          bool
	PrimaryLanguage struct {
		Name string
	}
}

type RepoInfoQuery struct {
	User struct {
		Repositories struct {
			PageInfo struct {
				HasNextPage bool
				EndCursor   string
			}
			Nodes []RepoInfo
		} `graphql:"repositories(first: 100, after: $cursor, orderBy: {field: PUSHED_AT, direction: DESC})"`
	} `graphql:"user(login: $login)"`
}

type RepoStats struct {
	StargazerCount int
	Languages      []*Lang
}

type Lang struct {
	Name  string
	Count int
}

func GetRepoStats(login string) (*RepoStats, error) {
	args := gin.H{
		"login":  githubv4.String(login),
		"cursor": (*githubv4.String)(nil),
	}

	stats := RepoStats{}
	langs := map[string]*Lang{}
	query := RepoInfoQuery{}

	// gather all repo stats
	for {
		err := gclient.Query(context.Background(), &query, args)
		if err != nil {
			return nil, errors.Wrap(err, "error fetching stargazers")
		}

		for _, node := range query.User.Repositories.Nodes {
			stats.StargazerCount += node.StargazerCount
			if !node.IsFork {
				// only count top languages for non-forks
				if langs[node.PrimaryLanguage.Name] == nil {
					langs[node.PrimaryLanguage.Name] = &Lang{node.PrimaryLanguage.Name, 1}
					stats.Languages = append(stats.Languages, langs[node.PrimaryLanguage.Name])
				} else {
					langs[node.PrimaryLanguage.Name].Count += 1
				}
			}
		}

		if !query.User.Repositories.PageInfo.HasNextPage {
			break
		}
		args["cursor"] = githubv4.NewString(githubv4.String(query.User.Repositories.PageInfo.EndCursor))
	}

	// descending sort
	sort.Slice(stats.Languages, func(i, j int) bool {
		return stats.Languages[i].Count > stats.Languages[j].Count
	})

	return &stats, nil
}
