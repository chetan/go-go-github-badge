package go_go_github_badge

import "testing"

func Test_isUserAllowed(t *testing.T) {
	SetAllowedUsers([]string{"chetan", "foobar"})
	type args struct {
		username string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{"chetan", args{"chetan"}, true},
		{"foobar", args{"foobar"}, true},
		{"baz", args{"baz"}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isUserAllowed(tt.args.username); got != tt.want {
				t.Errorf("isUserAllowed() = %v, want %v", got, tt.want)
			}
		})
	}
}
