package bj

type Participant interface {
	Name() string
}

type Action uint8

const (
	ActionHit Action = iota
	ActionStand
)

func RegisterAction(gs *GameState, participant Participant, action Action) {
	for _, p := range gs.Players {
		p.actor.OnAction(participant.Name(), action)
	}
}
