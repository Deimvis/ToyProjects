//go:generate stringer -type=Suit,Rank

package deck

import "fmt"

// Deck
// A Deck represents a deck of playing cards
type Deck []Card

// Draw takes number of cards to draw from the top of the deck
// and returns these cards.
// It returns an error if there are not enough cards in the deck
func (d *Deck) Draw(n int) ([]Card, error) {
	if len(*d) < n {
		return nil, fmt.Errorf("not enough cards in deck to draw: %d < %d", len(*d), n)
	}
	var cards []Card
	cards, *d = (*d)[:n], (*d)[n:]
	return cards, nil
}

// DrawOne draws one card from the top of the deck
// Basically, the same as Draw(1)
func (d *Deck) DrawOne() (Card, error) {
	cards, err := d.Draw(1)
	if err != nil {
		return Card{}, err
	}
	return cards[0], nil
}

// A Card represents playing card.
// The zero value for Card is not a valid card
type Card struct {
	Suit // may be Spade, Diamond, Club, Heart or Joker (special case)
	Rank // may be 2-10, Jack, Queen, King
}

// Represents Card as a string
func (c Card) String() string {
	if c.Suit == Joker {
		return c.Suit.String()
	}
	return fmt.Sprintf("%s of %ss", c.Rank.String(), c.Suit.String())
}

// NOTE: I actually disagree with using iota for such constants
// since there is a big chance for exporting cards with these constants somewhere else (e.g. in file or database)
// after that it will be impossible or very difficult to make any changes in the following code
// even though there's a small chance of any changes it's simply not worth to make it impossible to change
// nevertheless I used iota for consistency with gophercises
// More about an appropriate use of iota: https://www.gopherguides.com/articles/how-to-use-iota-in-golang

// Card's suit
type Suit uint8

const (
	Spade Suit = iota
	Diamond
	Club
	Heart
	Joker // special case
)

var suits = [...]Suit{Spade, Diamond, Club, Heart}

// Card's rank
type Rank uint8

const (
	Ace Rank = 1 + iota
	Two
	Three
	Four
	Five
	Six
	Seven
	Eight
	Nine
	Ten
	Jack
	Queen
	King
)

const (
	minRank = Ace
	maxRank = King
)

func absRank(c Card) int {
	return int(c.Suit)*int(maxRank) + int(c.Rank)
}
