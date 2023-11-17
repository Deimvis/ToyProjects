package bj

import (
	"errors"

	"github.com/Deimvis/toyprojects/deck_of_cards/deck"
)

type PlayerAPI struct {
	IsMyTurn       func() bool
	PlayerNames    func() []string
	ViewHand       func(playerName string) Hand
	ViewDealerCard func() deck.Card
	Hit            func()
	Stand          func()
}

func NewPlayerAPI(gs *GameState, p *Player) PlayerAPI {
	isMyTurn := func() bool {
		return gs.Stage == StagePlayerTurn && gs.CurrentPlayer().Name() == p.Name()
	}
	playerNames := func() []string {
		res := make([]string, len(gs.Players))
		for i := range gs.Players {
			res[i] = gs.Players[i].Name()
		}
		return res
	}
	viewHand := func(playerName string) Hand {
		for _, p := range gs.Players {
			if p.Name() == playerName {
				return p.hand
			}
		}
		return Hand{}
	}
	viewDealerCard := func() deck.Card {
		return gs.Dealer.FirstCard()
	}
	stand := func() {
		if !isMyTurn() {
			panic(notYourTurn{errors.New("it's not your turn anymore")})
		}
		RegisterAction(gs, p, ActionStand)
		if gs.curPlayer+1 < len(gs.Players) {
			gs.curPlayer++
		} else {
			gs.Stage = StageDealerTurn
		}
	}
	hit := func() {
		if !isMyTurn() {
			panic(notYourTurn{errors.New("it's not your turn anymore")})
		}
		RegisterAction(gs, p, ActionHit)
		card, _ := gs.Deck.DrawOne()
		p.Take(card)
		if p.hand.Score() > 21 {
			stand()
		}
	}
	api := PlayerAPI{
		IsMyTurn: func() bool {
			return isMyTurn()
		},
		PlayerNames: func() []string {
			return playerNames()
		},
		ViewHand: func(playerName string) Hand {
			return viewHand(playerName)
		},
		ViewDealerCard: func() deck.Card {
			return viewDealerCard()
		},
		Hit: func() {
			hit()
		},
		Stand: func() {
			stand()
		},
	}
	return api
}
