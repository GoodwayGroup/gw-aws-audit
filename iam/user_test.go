package iam

import (
	"fmt"
	"github.com/logrusorgru/aurora/v3"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestIamUser_ARN(t *testing.T) {
	is := assert.New(t)

	t.Run("should return the literal string", func(t *testing.T) {
		test := "test"
		u := User{
			arn: &test,
		}

		is.Equal(test, u.ARN())
	})

	t.Run("should return empty string when there is no value", func(t *testing.T) {
		u := User{}

		is.Equal("", u.ARN())
	})
}

func TestIamUser_UserName(t *testing.T) {
	is := assert.New(t)

	t.Run("should return the literal string", func(t *testing.T) {
		test := "test"
		u := User{
			userName: &test,
		}

		is.Equal(test, u.UserName())
	})

	t.Run("should return empty string when there is no value", func(t *testing.T) {
		u := User{}

		is.Equal("", u.UserName())
	})
}

func TestIamUser_ID(t *testing.T) {
	is := assert.New(t)

	t.Run("should return the literal string", func(t *testing.T) {
		test := "test"
		u := User{
			userID: &test,
		}

		is.Equal(test, u.ID())
	})

	t.Run("should return empty string when there is no value", func(t *testing.T) {
		u := User{}

		is.Equal("", u.ID())
	})
}

func TestIamUser_HasConsoleAccess(t *testing.T) {
	is := assert.New(t)

	t.Run("should return the value set", func(t *testing.T) {
		u := User{
			hasConsoleAccess: true,
		}

		is.True(u.HasConsoleAccess())
	})

	t.Run("should return false by default", func(t *testing.T) {
		u := User{}

		is.False(u.HasConsoleAccess())
	})
}

func TestIamUser_LastLogin(t *testing.T) {
	is := assert.New(t)

	t.Run("should return the value set", func(t *testing.T) {
		test := time.Now()
		u := User{
			passwordLastUsed: &test,
		}

		is.Equal(test, u.LastLogin())
	})

	t.Run("should return zeroed time.Time", func(t *testing.T) {
		u := User{}

		is.Equal(time.Time{}, u.LastLogin())
	})
}

func TestIamUser_LastLoginDuration(t *testing.T) {
	is := assert.New(t)

	t.Run("should return the value set", func(t *testing.T) {
		test := time.Date(2020, time.January, 1, 0, 0, 0, 0, time.UTC)
		u := User{
			passwordLastUsed: &test,
		}

		is.Regexp(`^\d+\s\w+$`, u.LastLoginDuration())
	})

	t.Run("should return zeroed time.Time", func(t *testing.T) {
		u := User{}

		is.Equal("", u.LastLoginDuration())
	})
}

func TestIamUser_CreatedDate(t *testing.T) {
	is := assert.New(t)

	t.Run("should return the value set", func(t *testing.T) {
		test := time.Now()
		u := User{
			createDate: &test,
		}

		is.Equal(test, u.CreatedDate())
	})

	t.Run("should return zeroed time.Time", func(t *testing.T) {
		u := User{}

		is.Equal(time.Time{}, u.CreatedDate())
	})
}

func TestIamUser_CreatedDateDuration(t *testing.T) {
	is := assert.New(t)

	t.Run("should return the value set", func(t *testing.T) {
		test := time.Date(2020, time.January, 1, 0, 0, 0, 0, time.UTC)
		u := User{
			createDate: &test,
		}

		is.Regexp(`^\d+\s\w+$`, u.CreatedDateDuration())
	})

	t.Run("should return zeroed time.Time", func(t *testing.T) {
		u := User{}

		is.Equal("", u.CreatedDateDuration())
	})
}

func TestIamUser_HasAccessKeys(t *testing.T) {
	is := assert.New(t)

	t.Run("should return the value set", func(t *testing.T) {
		u := User{
			accessKeys: []*AccessKey{{}},
		}

		is.True(u.HasAccessKeys())
	})

	t.Run("should return false by default", func(t *testing.T) {
		u := User{}

		is.False(u.HasAccessKeys())
	})
}

func TestIamUser_AccessKeysCount(t *testing.T) {
	is := assert.New(t)

	t.Run("should return the value set", func(t *testing.T) {
		u := User{
			accessKeys: []*AccessKey{{}, {}},
		}

		is.Equal(2, u.AccessKeysCount())
	})

	t.Run("should return false by default", func(t *testing.T) {
		u := User{}

		is.Equal(0, u.AccessKeysCount())
	})
}

func TestIamUser_CheckStatus(t *testing.T) {
	is := assert.New(t)

	var tests = []struct {
		input User
		want  string
	}{
		{
			input: User{},
			want:  "pass",
		},
		{
			input: User{
				hasConsoleAccess: true,
			},
			want: "fail",
		},
		{
			input: User{
				hasConsoleAccess: false,
				accessKeys:       []*AccessKey{{}},
			},
			want: "warn",
		},
	}

	for _, tt := range tests {
		testname := fmt.Sprintf("should return status [%s]", tt.want)
		t.Run(testname, func(t *testing.T) {
			is.Equal(tt.want, tt.input.CheckStatus())
		})
	}
}

func TestIamUser_FormattedCheckStatus(t *testing.T) {
	is := assert.New(t)

	var tests = []struct {
		input User
		want  string
	}{
		{
			input: User{},
			want:  aurora.Green("PASS").String(),
		},
		{
			input: User{
				hasConsoleAccess: true,
			},
			want: aurora.Red("FAIL").String(),
		},
		{
			input: User{
				hasConsoleAccess: false,
				accessKeys:       []*AccessKey{{}},
			},
			want: aurora.Yellow("WARN").String(),
		},
	}

	for _, tt := range tests {
		testName := fmt.Sprintf("should return status [%s]", tt.want)
		t.Run(testName, func(t *testing.T) {
			is.Equal(tt.want, tt.input.FormattedCheckStatus())
		})
	}
}

func TestIamUser_FormattedLastLoginDateDuration(t *testing.T) {
	is := assert.New(t)

	t.Run("should return the value set", func(t *testing.T) {
		test := time.Date(2020, time.January, 1, 0, 0, 0, 0, time.UTC)
		u := User{
			hasConsoleAccess: true,
			passwordLastUsed: &test,
		}

		is.Regexp(`^\d+\s\w+$`, u.FormattedLastLoginDateDuration())
	})

	t.Run("should return zeroed time.Time", func(t *testing.T) {
		u := User{}

		is.Equal(aurora.Gray(8, "NONE").String(), u.FormattedLastLoginDateDuration())
	})
}
