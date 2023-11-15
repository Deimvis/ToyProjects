package deck

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func ExampleCard() {
	fmt.Println(Card{Rank: Ace, Suit: Heart})
	fmt.Println(Card{Rank: Two, Suit: Spade})
	fmt.Println(Card{Rank: Nine, Suit: Diamond})
	fmt.Println(Card{Rank: Jack, Suit: Club})
	fmt.Println(Card{Suit: Joker})

	// Output:
	// Ace of Hearts
	// Two of Spades
	// Nine of Diamonds
	// Jack of Clubs
	// Joker
}

func TestCard_String(t *testing.T) {
	testCases := []struct {
		card     Card
		expected string
	}{
		{Card{Suit: Joker}, "Joker"},
	}
	for _, tc := range testCases {
		require.Equal(t, tc.expected, tc.card.String())
	}
}
