package main

import (
	"flag"
	"fmt"

	bj "github.com/Deimvis/toyprojects/blackjack/src"
)

func main() {
	name := flag.String("name", "Player", "Player name")
	flag.Parse()
	gs := bj.NewGame([]bj.Player{
		bj.NewPlayer(*name, bj.STDINPlayerActor{}),
	})
	gs.Run()
	EndGame(*name, &gs)
}

func EndGame(name string, gs *bj.GameState) {
	fmt.Println("==FINAL HANDS==")
	for _, p := range gs.Players {
		bj.PrintInfo(p.Name(), p.Hand())
	}
	bj.PrintInfo("Dealer", gs.Dealer.Hand())
	pScore := gs.GetPlayer(name).Hand().Score()
	dScore := gs.Dealer.Hand().Score()
	switch {
	case pScore > 21:
		fmt.Println("You busted")
	case pScore > 21:
		fmt.Println("Dealer busted")
	case pScore > dScore:
		fmt.Println("You win!")
	case pScore < dScore:
		fmt.Println("You lost")
	case pScore == dScore:
		fmt.Println("Draw")
	}
}
