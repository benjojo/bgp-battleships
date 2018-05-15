package main

import (
	"flag"
	"log"
)

func main() {
	flag.Parse()
	log.Printf("yup")

	LocalB := battleShipBoard{}
	RemoteB := battleShipBoard{}

	N.Board[1][5] = stateShip
	N.Board[1][2] = stateHit
	N.Board[2][2] = stateAttempt

	N.Draw()
}
