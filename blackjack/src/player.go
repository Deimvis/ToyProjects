package bj

import (
	"fmt"
	"sort"

	"github.com/Deimvis/toyprojects/deck_of_cards/deck"
)

func NewPlayer(name string, actor PlayerActor) Player {
	p := Player{
		name:  name,
		hand:  make(Hand, 0, 5),
		actor: actor,
	}
	return p
}

type Player struct {
	name  string
	hand  Hand
	actor PlayerActor
}

func (p Player) Name() string {
	return p.name
}

func (p Player) Hand() Hand {
	return p.hand
}

func (p *Player) Take(c deck.Card) {
	p.hand.Add(c)
}

func (p Player) Play(api PlayerAPI) {
	p.actor.Play(api)
}

type PlayerActor interface {
	Play(api PlayerAPI)
	OnAction(participantName string, action Action)
}

type STDINPlayerActor struct {
}

func (a STDINPlayerActor) Play(api PlayerAPI) {
	for api.IsMyTurn() {
		names := api.PlayerNames()
		sort.Slice(names, func(i, j int) bool {
			return api.ViewHand(names[i]).Score() > api.ViewHand(names[j]).Score()
		})
		for _, name := range names {
			hand := api.ViewHand(name)
			PrintInfo(name, hand)
		}
		PrintInfo("Dealer", Hand{api.ViewDealerCard()})
		fmt.Println("What will you do? (h)it, (s)tand")
		var input string
		fmt.Scanf("%s\n", &input)
		switch input {
		case "h":
			api.Hit()
		case "s":
			api.Stand()
		}
	}
}

func (a STDINPlayerActor) OnAction(participantName string, action Action) {
	var actionName string
	switch action {
	case ActionHit:
		actionName = "Hit!"
	case ActionStand:
		actionName = "Stand"
	}
	fmt.Printf("> %s: %s\n", participantName, actionName)
}
