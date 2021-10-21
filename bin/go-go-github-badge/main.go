package main

import (
	"os"
	"strings"

	go_go_github_badge "github.com/chetan/go-go-github-badge"
)

func main() {
	allowedUsers := os.Getenv("ALLOWED_USERS")
	var allowed []string
	if allowedUsers != "" {
		allowed = strings.Split(allowedUsers, ",")
	}

	go_go_github_badge.SetAllowedUsers(allowed)
	go_go_github_badge.Run()
}
