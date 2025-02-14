package main

import (
	"fmt"
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

	if len(os.Args) >= 3 && os.Args[1] == "gen" {
		go_go_github_badge.CreateClient()
		badge, err := go_go_github_badge.Generate(os.Args[2])
		if err != nil {
			fmt.Println("Error:", err)
			os.Exit(1)
		}
		fmt.Println(badge)
		return
	}

	go_go_github_badge.Run()
}
