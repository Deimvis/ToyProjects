package bj

type DealerAPI struct {
	IsMyTurn func() bool
	Hand     func() Hand
	Hit      func()
	Stand    func()
}

func NewDealerAPI(gs *GameState, d *Dealer) DealerAPI {
	isMyTurn := func() bool {
		return gs.Stage == StageDealerTurn
	}
	hand := func() Hand {
		return gs.Dealer.Hand()
	}
	stand := func() {
		if !isMyTurn() {
			panic("It's not your turn anymore")
		}
		RegisterAction(gs, d, ActionStand)
		gs.Stage = StageHandOver
	}
	hit := func() {
		if !isMyTurn() {
			panic("It's not your turn anymore")
		}
		RegisterAction(gs, d, ActionHit)
		card, _ := gs.Deck.DrawOne()
		d.Take(card)
		if d.hand.Score() > 21 {
			stand()
		}
	}
	api := DealerAPI{
		IsMyTurn: func() bool {
			return isMyTurn()
		},
		Hand: func() Hand {
			return hand()
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
