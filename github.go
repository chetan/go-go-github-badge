package go_go_github_badge

import (
	"context"
	"os"
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

func GetLatestContributions(user *github.User) (*LatestContributions, error) {
	to := time.Now()
	from := to.Add(-time.Hour * 24 * 90)
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

type StargazerCountRepo struct {
	StargazerCount int
}

type StargazerCountQuery struct {
	User struct {
		Repositories struct {
			PageInfo struct {
				HasNextPage bool
				EndCursor   string
			}
			Nodes []StargazerCountRepo
		} `graphql:"repositories(first: 100, after: $cursor, orderBy: {field: PUSHED_AT, direction: DESC})"`
	} `graphql:"user(login: $login)"`
}

func GetStargazerCount(login string) (int, error) {
	args := gin.H{
		"login":  githubv4.String(login),
		"cursor": (*githubv4.String)(nil),
	}

	totalCount := 0
	query := StargazerCountQuery{}

	for {
		err := gclient.Query(context.Background(), &query, args)
		if err != nil {
			return 0, errors.Wrap(err, "error fetching stargazers")
		}

		for _, node := range query.User.Repositories.Nodes {
			totalCount += node.StargazerCount
		}

		if !query.User.Repositories.PageInfo.HasNextPage {
			break
		}
		args["cursor"] = githubv4.NewString(githubv4.String(query.User.Repositories.PageInfo.EndCursor))
	}

	return totalCount, nil
}
