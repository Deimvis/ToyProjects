package bj

import (
	"github.com/Deimvis/toyprojects/deck_of_cards/deck"
	"github.com/barkimedes/go-deepcopy"
)

func NewGame(players []Player) GameState {
	resetMaxNameWidth(players)
	return GameState{Deck: deck.New(deck.Copies(3), deck.Shuffle), Players: players, Dealer: NewDealer()}
}

type GameState struct {
	Deck deck.Deck
	Stage
	Players []Player
	Dealer  Dealer

	curPlayer int // makes sense only while StagePlayerTurn
}

func (gs *GameState) Run() {
	gs.processDealStage()
	gs.processPlayStage()
	gs.processDealerPlayStage()
}

func (gs *GameState) CurrentPlayer() Player {
	if gs.Stage != StagePlayerTurn {
		panic("It's not a player turn")
	}
	return gs.Players[gs.curPlayer]
}

func (gs GameState) GetPlayer(name string) Player {
	for _, p := range gs.Players {
		if p.Name() == name {
			return p
		}
	}
	return Player{}
}

func (gs GameState) Clone() GameState {
	return *(deepcopy.MustAnything(&gs).(*GameState))
}

type Stage uint8

const (
	StagePlayerTurn Stage = iota
	StageDealerTurn
	StageHandOver
)
