package bj

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRestMaxWidth(t *testing.T) {
	testCases := []struct {
		players  []Player
		expected int
	}{
		{[]Player{NewPlayer("", STDINPlayerActor{})}, len("Dealer")},
		{[]Player{NewPlayer("len10     ", STDINPlayerActor{})}, 10},
		{[]Player{NewPlayer("len10     ", STDINPlayerActor{}), NewPlayer("len11      ", STDINPlayerActor{})}, 11},
	}
	for _, tc := range testCases {
		resetMaxNameWidth(tc.players)
		require.Equal(t, tc.expected, maxWidth)
	}
}
