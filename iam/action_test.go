package iam

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_determineCheck(t *testing.T) {
	is := assert.New(t)
	tests := []struct {
		threshold int64
		units     string
		want      int64
	}{
		{threshold: 1, units: "hours", want: 1},
		{threshold: 1, units: "days", want: 24},
		{threshold: 1, units: "weeks", want: 168},
		{threshold: 1, units: "months", want: 720},
		{threshold: 1, units: "", want: 24},
	}

	for _, tc := range tests {
		is.Equal(tc.want, determineCheck(tc.threshold, tc.units))
	}
}
