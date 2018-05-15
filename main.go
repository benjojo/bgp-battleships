package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
)

func main() {
	flag.Parse()
	log.Printf("yup")

	LocalB := battleShipBoard{}
	// RemoteB := battleShipBoard{}

	LocalB.Board[1][5] = stateShip
	LocalB.Board[1][2] = stateHit
	LocalB.Board[2][2] = stateAttempt

	LocalB.Draw()

	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("Next Move> ")
		text, _ := reader.ReadString('\n')
		if len(text) != 3 {
			log.Printf("wrong length of command %d", len(text))
			continue
		}
		x, y := cordsToNumbers(text)
		if x == -1 || y == -1 {
			continue
		}
		LocalB.Board[y][x] = stateHit
		LocalB.Draw()
	}

}
