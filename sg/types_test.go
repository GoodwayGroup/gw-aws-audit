package sg

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSecurityGroup_ID(t *testing.T) {
	is := assert.New(t)

	t.Run("should return the literal string", func(t *testing.T) {
		test := "sg-1234"
		sg := SecurityGroup{
			id: test,
		}

		is.Equal(test, sg.ID())
	})

	t.Run("should return empty string when there is no value", func(t *testing.T) {
		sg := SecurityGroup{}

		is.Equal("", sg.ID())
	})
}

func TestSecurityGroup_Name(t *testing.T) {
	is := assert.New(t)

	t.Run("should return the literal string", func(t *testing.T) {
		test := "test-name"
		sg := SecurityGroup{
			name: test,
		}

		is.Equal(test, sg.Name())
	})

	t.Run("should return empty string when there is no value", func(t *testing.T) {
		sg := SecurityGroup{}

		is.Equal("", sg.Name())
	})
}

func TestSecurityGroup_Attachments(t *testing.T) {
	is := assert.New(t)

	t.Run("should return the literal string", func(t *testing.T) {
		attachments := map[string]int{"ec2": 1, "rds": 2}
		sg := SecurityGroup{
			attached: attachments,
		}

		is.Equal(attachments, sg.Attachments())
	})

	t.Run("should return empty string when there is no value", func(t *testing.T) {
		sg := SecurityGroup{}

		var test map[string]int
		is.Equal(test, sg.Attachments())
	})
}

func TestSecurityGroup_GetAttachmentsAsString(t *testing.T) {
	is := assert.New(t)

	t.Run("should return the literal string", func(t *testing.T) {
		attachments := map[string]int{"ec2": 1, "rds": 2}
		sg := SecurityGroup{
			attached: attachments,
		}

		is.Equal("ec2: 1 rds: 2", sg.GetAttachmentsAsString())
	})

	t.Run("should return empty string when there is no value", func(t *testing.T) {
		sg := SecurityGroup{}

		is.Equal("", sg.GetAttachmentsAsString())
	})
}
