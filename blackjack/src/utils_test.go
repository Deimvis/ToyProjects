package bj

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLjust(t *testing.T) {
	testCases := []struct {
		s        string
		n        int
		fill     string
		expected string
	}{
		{"Hello, World!", 3, " ", "Hello, World!"},
		{"Hello, World!", 0, " ", "Hello, World!"},
		{"Hello, World!", -100, " ", "Hello, World!"},
		{"Hello, World!", 15, " ", "Hello, World!  "},
		{"", 5, "-", "-----"},
		{"aaa", 5, "-", "aaa--"},
	}
	for _, tc := range testCases {
		require.Equal(t, tc.expected, ljust(tc.s, tc.n, tc.fill))
	}
}

func TsetLjustPanics(t *testing.T) {
	require.Panics(t, func() { ljust("", 0, "some text here") })
	require.Panics(t, func() { ljust("", 0, "cant even have length = two") })
	require.Panics(t, func() { ljust("", 0, "22") })
	require.Panics(t, func() { ljust("", 0, "and length = 0") })
	require.Panics(t, func() { ljust("", 0, "") })
}
