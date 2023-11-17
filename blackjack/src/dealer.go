package bj

import "github.com/Deimvis/toyprojects/deck_of_cards/deck"

func NewDealer() Dealer {
	return Dealer{hand: make(Hand, 0, 5)}
}

type Dealer struct {
	hand Hand
}

func (d Dealer) Name() string {
	return "Dealer"
}

func (d Dealer) Hand() Hand {
	return d.hand
}

func (d Dealer) FirstCard() deck.Card {
	return d.hand[0]
}

func (d *Dealer) Take(c deck.Card) {
	d.hand.Add(c)
}

func (d Dealer) Play(api DealerAPI) {
	for api.IsMyTurn() {
		if api.Hand().Score() <= 16 || (api.Hand().Score() == 17 && api.Hand().MinScore() != 17) {
			api.Hit()
		} else {
			api.Stand()
		}
	}
}
