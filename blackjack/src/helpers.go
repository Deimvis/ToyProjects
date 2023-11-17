package bj

import "fmt"

var defaultMaxWidth int = len("Dealer")
var maxWidth int = defaultMaxWidth

func resetMaxNameWidth(players []Player) {
	maxWidth = defaultMaxWidth
	for _, p := range players {
		maxWidth = max(maxWidth, len(p.Name()))
	}
}

func PrintInfo(name string, hand Hand) {
	fmt.Printf("%s: %d (%s)\n", ljust(name, maxWidth+1, " "), hand.Score(), hand.String())
}
