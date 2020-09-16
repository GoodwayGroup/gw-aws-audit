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
		u := iamUser{
			arn: &test,
		}

		is.Equal(test, u.ARN())
	})

	t.Run("should return empty string when there is no value", func(t *testing.T) {
		u := iamUser{}

		is.Equal("", u.ARN())
	})
}

func TestIamUser_UserName(t *testing.T) {
	is := assert.New(t)

	t.Run("should return the literal string", func(t *testing.T) {
		test := "test"
		u := iamUser{
			userName: &test,
		}

		is.Equal(test, u.UserName())
	})

	t.Run("should return empty string when there is no value", func(t *testing.T) {
		u := iamUser{}

		is.Equal("", u.UserName())
	})
}

func TestIamUser_ID(t *testing.T) {
	is := assert.New(t)

	t.Run("should return the literal string", func(t *testing.T) {
		test := "test"
		u := iamUser{
			userID: &test,
		}

		is.Equal(test, u.ID())
	})

	t.Run("should return empty string when there is no value", func(t *testing.T) {
		u := iamUser{}

		is.Equal("", u.ID())
	})
}

func TestIamUser_HasConsoleAccess(t *testing.T) {
	is := assert.New(t)

	t.Run("should return the value set", func(t *testing.T) {
		u := iamUser{
			hasConsoleAccess: true,
		}

		is.True(u.HasConsoleAccess())
	})

	t.Run("should return false by default", func(t *testing.T) {
		u := iamUser{}

		is.False(u.HasConsoleAccess())
	})
}

func TestIamUser_LastLogin(t *testing.T) {
	is := assert.New(t)

	t.Run("should return the value set", func(t *testing.T) {
		test := time.Now()
		u := iamUser{
			passwordLastUsed: &test,
		}

		is.Equal(test, u.LastLogin())
	})

	t.Run("should return zeroed time.Time", func(t *testing.T) {
		u := iamUser{}

		is.Equal(time.Time{}, u.LastLogin())
	})
}

func TestIamUser_LastLoginDuration(t *testing.T) {
	is := assert.New(t)

	t.Run("should return the value set", func(t *testing.T) {
		test := time.Date(2020, time.January, 1, 0, 0, 0, 0, time.UTC)
		u := iamUser{
			passwordLastUsed: &test,
		}

		is.Regexp(`^\d+\s\w+$`, u.LastLoginDuration())
	})

	t.Run("should return zeroed time.Time", func(t *testing.T) {
		u := iamUser{}

		is.Equal("", u.LastLoginDuration())
	})
}

func TestIamUser_CreatedDate(t *testing.T) {
	is := assert.New(t)

	t.Run("should return the value set", func(t *testing.T) {
		test := time.Now()
		u := iamUser{
			createDate: &test,
		}

		is.Equal(test, u.CreatedDate())
	})

	t.Run("should return zeroed time.Time", func(t *testing.T) {
		u := iamUser{}

		is.Equal(time.Time{}, u.CreatedDate())
	})
}

func TestIamUser_CreatedDateDuration(t *testing.T) {
	is := assert.New(t)

	t.Run("should return the value set", func(t *testing.T) {
		test := time.Date(2020, time.January, 1, 0, 0, 0, 0, time.UTC)
		u := iamUser{
			createDate: &test,
		}

		is.Regexp(`^\d+\s\w+$`, u.CreatedDateDuration())
	})

	t.Run("should return zeroed time.Time", func(t *testing.T) {
		u := iamUser{}

		is.Equal("", u.CreatedDateDuration())
	})
}

func TestIamUser_HasAccessKeys(t *testing.T) {
	is := assert.New(t)

	t.Run("should return the value set", func(t *testing.T) {
		u := iamUser{
			accessKeys: []*accessKey{{}},
		}

		is.True(u.HasAccessKeys())
	})

	t.Run("should return false by default", func(t *testing.T) {
		u := iamUser{}

		is.False(u.HasAccessKeys())
	})
}

func TestIamUser_AccessKeysCount(t *testing.T) {
	is := assert.New(t)

	t.Run("should return the value set", func(t *testing.T) {
		u := iamUser{
			accessKeys: []*accessKey{{}, {}},
		}

		is.Equal(2, u.AccessKeysCount())
	})

	t.Run("should return false by default", func(t *testing.T) {
		u := iamUser{}

		is.Equal(0, u.AccessKeysCount())
	})
}

func TestIamUser_CheckStatus(t *testing.T) {
	is := assert.New(t)

	var tests = []struct {
		input iamUser
		want  string
	}{
		{
			input: iamUser{},
			want:  "pass",
		},
		{
			input: iamUser{
				hasConsoleAccess: true,
			},
			want: "fail",
		},
		{
			input: iamUser{
				hasConsoleAccess: false,
				accessKeys:       []*accessKey{{}},
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
		input iamUser
		want  string
	}{
		{
			input: iamUser{},
			want:  aurora.Green("PASS").String(),
		},
		{
			input: iamUser{
				hasConsoleAccess: true,
			},
			want: aurora.Red("FAIL").String(),
		},
		{
			input: iamUser{
				hasConsoleAccess: false,
				accessKeys:       []*accessKey{{}},
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
		u := iamUser{
			hasConsoleAccess: true,
			passwordLastUsed: &test,
		}

		is.Regexp(`^\d+\s\w+$`, u.FormattedLastLoginDateDuration())
	})

	t.Run("should return zeroed time.Time", func(t *testing.T) {
		u := iamUser{}

		is.Equal(aurora.Gray(8, "NONE").String(), u.FormattedLastLoginDateDuration())
	})
}
