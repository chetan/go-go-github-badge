package go_go_github_badge

import (
	"reflect"
	"testing"

	"github.com/google/go-github/v39/github"
)

func init() {
}

func TestGetUser(t *testing.T) {
	type args struct {
		username string
	}
	tests := []struct {
		name    string
		args    args
		want    *github.User
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetUser(tt.args.username)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetUser() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetLatestContributions(t *testing.T) {
	CreateClient()

	type args struct {
		user *github.User
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
		{"chetan", args{&github.User{Login: ptr("chetan")}}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := GetLatestContributions(tt.args.user); (err != nil) != tt.wantErr {
				t.Errorf("GetLatestContributions() error = %v, wantErr %v", err, tt.wantErr)
			}
			// t.FailNow()
		})
	}
}

func ptr(s string) *string {
	return &s
}
