package deck

import (
	"math/rand"
	"testing"

	"github.com/stretchr/testify/require"
)

const expectedDeckSz = 13 * 4

func TestNew(t *testing.T) {
	cards := New()
	require.Equal(t, expectedDeckSz, len(cards))
}

func TestDefaultSort(t *testing.T) {
	cards := New(DefaultSort)
	require.Equal(t, Card{Rank: Ace, Suit: Spade}, cards[0])
}

func TestSort(t *testing.T) {
	greater := func(cards []Card) func(i, j int) bool {
		return func(i, j int) bool {
			return absRank(cards[i]) > absRank(cards[j])
		}
	}
	cards := New(Sort(greater))
	require.Equal(t, Card{Rank: Ace, Suit: Spade}, cards[len(cards)-1])
}

func TestJokers(t *testing.T) {
	cards := New(Jokers(3))
	count := 0
	for _, c := range cards {
		if c.Suit == Joker {
			count++
		}
	}
	require.Equal(t, 3, count)
}

func TestFilter(t *testing.T) {
	filter := func(card Card) bool {
		return card.Rank != Two && card.Rank != Three
	}
	cards := New(Filter(filter))
	for _, c := range cards {
		require.NotEqual(t, Two, c.Rank)
		require.NotEqual(t, Three, c.Rank)
	}
}

func TestExclude(t *testing.T) {
	exclude := func(card Card) bool {
		return card.Rank != Two && card.Rank != Three
	}
	cards := New(Filter(exclude))
	for _, c := range cards {
		require.Contains(t, []Rank{Two, Three}, c.Rank)
	}
}

func TestCopies(t *testing.T) {
	cards := New(Copies(3))
	require.Equal(t, expectedDeckSz*3, len(cards))
}

func TestShuffle(t *testing.T) {
	shuffleRand = rand.New(rand.NewSource(1))
	cards := New(Shuffle)
	require.Equal(t, expectedShuffleResult, cards)
}

var expectedShuffleResult = []Card{
	{Club, King},
	{Club, Ace},
	{Heart, Five},
	{Heart, King},
	{Diamond, Ace},
	{Diamond, Ten},
	{Heart, Eight},
	{Club, Queen},
	{Diamond, Six},
	{Club, Nine},
	{Diamond, Queen},
	{Spade, Three},
	{Club, Three},
	{Diamond, Three},
	{Heart, Two},
	{Spade, Ten},
	{Spade, Two},
	{Club, Four},
	{Heart, Three},
	{Heart, Seven},
	{Spade, Ace},
	{Club, Five},
	{Diamond, Four},
	{Heart, Queen},
	{Spade, Six},
	{Diamond, Seven},
	{Club, Two},
	{Heart, Ace},
	{Diamond, King},
	{Club, Jack},
	{Spade, Queen},
	{Spade, Seven},
	{Heart, Six},
	{Diamond, Jack},
	{Club, Ten},
	{Spade, Jack},
	{Diamond, Five},
	{Heart, Four},
	{Diamond, Two},
	{Spade, Nine},
	{Heart, Jack},
	{Heart, Ten},
	{Spade, King},
	{Spade, Five},
	{Spade, Eight},
	{Spade, Four},
	{Club, Seven},
	{Diamond, Eight},
	{Diamond, Nine},
	{Club, Eight},
	{Heart, Nine},
	{Club, Six},
}
