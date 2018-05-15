package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"time"
)

func main() {
	startfirst := flag.Bool("startfirst", false, "set this if you are starting first")
	flag.Parse()
	log.Printf("yup")

	LocalB := makeBoard()
	RemoteB := battleShipBoard{}

	LocalB.Draw()

	gameCounter := 0
	hitmiss := 0

	fmt.Print("Your Side                   Player Two\n")
	fmt.Print(combineBoard(LocalB, RemoteB))

	reader := bufio.NewReader(os.Stdin)
	for {
		var text string
		var x, y int
		if *startfirst {
			fmt.Printf("[%06d] Next Move> ", gameCounter)
			text, _ = reader.ReadString('\n')
			if len(text) != 3 {
				log.Printf("wrong length of command %d", len(text))
				continue
			}
			x, y = cordsToNumbers(text)
			if x == -1 || y == -1 {
				continue
			}

			fmt.Printf("Firing on %s...", text)
			writeBGP(gameCounter, x, y, hitmiss)
		}

		*startfirst = true

		fmt.Printf("waiting on players responce...\n")

		for {
			time.Sleep(time.Second)
			var err error
			var nx, ny int
			tempgameCounter := 0
			tempgameCounter, nx, ny, hitmiss, err = readBGP()
			if err != nil {
				fmt.Print(".")
				continue
			}
			fmt.Print("!")

			if tempgameCounter > gameCounter || gameCounter == 0 {
				gameCounter = tempgameCounter + 1
				// !! New move has happened

				// First, process if we got a hit or not.
				if hitmiss == 1 {
					RemoteB.Board[y][x] = stateHit
					log.Printf("It's a Hit!")
				} else {
					RemoteB.Board[y][x] = stateAttempt
					log.Printf("It's a Miss!")
				}

				if !(tempgameCounter == 1 && *startfirst) {
					// Now... did we get hit?
					if LocalB.Board[ny][nx] == stateShip {
						hitmiss = 1
						LocalB.Board[ny][nx] = stateHit
					} else {
						hitmiss = 0
						LocalB.Board[ny][nx] = stateAttempt
					}
				}
				break
			}
		}

		fmt.Print("Your Side                   Player Two\n")
		fmt.Print(combineBoard(LocalB, RemoteB))
	}

}
