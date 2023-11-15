package deck

import (
	"math/rand"
	"sort"
	"time"
)

// New returns a new deck (slice of cards)
func New(opts ...func([]Card) []Card) []Card {
	var cards []Card
	for _, suit := range suits {
		for rank := minRank; rank <= maxRank; rank++ {
			cards = append(cards, Card{Suit: suit, Rank: rank})
		}
	}
	for _, opt := range opts {
		cards = opt(cards)
	}
	return cards
}

// DefaultSort is an option for New function that sorts deck in some preserved order
func DefaultSort(cards []Card) []Card {
	return Sort(Less)(cards)
}

// Sort takes cards comparator as an argument and returns an option for New function
func Sort(less func(cards []Card) func(i, j int) bool) func([]Card) []Card {
	return func(cards []Card) []Card {
		sort.Slice(cards, less(cards))
		return cards
	}
}

// Less takes cards as an argument and returns cards comparator.
// The card order of returned comparator is defined as follows:
// < Spades, Diamonds, Clubs, Hearts, Joker
// tie is resolved by rank as follows:
// < Ace, Two, Three, ..., Ten, Jack, Queen, King
func Less(cards []Card) func(i, j int) bool {
	return func(i, j int) bool {
		return absRank(cards[i]) < absRank(cards[j])
	}
}

var shuffleRand = rand.New(rand.NewSource(time.Now().Unix()))

// Shuffle is an option for New function that shuffles given deck
func Shuffle(cards []Card) []Card {
	shuffleRand.Shuffle(len(cards), func(i, j int) { cards[i], cards[j] = cards[j], cards[i] })
	return cards
}

// Jokers takes number of jokers to be inserted into deck and
// returns an option for New function
func Jokers(n int) func([]Card) []Card {
	return func(cards []Card) []Card {
		for i := 0; i < n; i++ {
			cards = append(cards, Card{
				Rank: Rank(i),
				Suit: Joker,
			})
		}
		return cards
	}
}

// Exclude takes function that reports whether card should be excluded from deck
// It returns an option for New function
func Exclude(f func(card Card) bool) func([]Card) []Card {
	return Filter(func(card Card) bool {
		return !f(card)
	})
}

// Filter takes function that reports whether card should be kept in deck
// It returns an option for New function
func Filter(f func(card Card) bool) func([]Card) []Card {
	return func(cards []Card) []Card {
		sz := 0
		for i := range cards {
			if f(cards[i]) {
				cards[sz] = cards[i]
			}
		}
		cards = cards[:sz]
		return cards
	}
}

// Decks takes a number of copies of decks to be created
// and returns an option for New funciton.
// For example: Decks(3) means that deck should contain 3 copies of itself
func Decks(n int) func([]Card) []Card {
	return func(cards []Card) []Card {
		return repeat(cards, n)
	}
}
