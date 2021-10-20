package go_go_github_badge

import (
	"reflect"
	"testing"

	"github.com/google/go-github/v39/github"
	"github.com/joho/godotenv"
)

func init() {
	err := godotenv.Overload()
	if err != nil {
		panic(err)
	}
	CreateClient()
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
	type args struct {
		user *github.User
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
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
