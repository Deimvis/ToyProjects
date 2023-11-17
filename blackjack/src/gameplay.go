package bj

func (gs *GameState) processDealStage() {
	for i := 0; i < 2; i++ {
		for pInd := range gs.Players {
			card, _ := gs.Deck.DrawOne()
			gs.Players[pInd].Take(card)
		}
		card, _ := gs.Deck.DrawOne()
		(&gs.Dealer).Take(card)
	}
	gs.Stage = StagePlayerTurn
}

func (gs *GameState) processPlayStage() {
	for i, p := range gs.Players {
		defer func() {
			if r := recover(); r != nil {
				if _, ok := r.(notYourTurn); !ok {
					panic(r)
				}
			}
		}()
		p.Play(NewPlayerAPI(gs, &gs.Players[i]))
	}
	gs.Stage = StageDealerTurn
}

func (gs *GameState) processDealerPlayStage() {
	gs.Dealer.Play(NewDealerAPI(gs, &gs.Dealer))
	gs.Stage = StageHandOver
}
