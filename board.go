package main

import (
	cr "crypto/rand"
	"fmt"
	"math"
	"math/big"
	"math/rand"
	"strconv"
	"strings"

	"github.com/mgutz/ansi"
)

type boardState int

const square string = "■"

const (
	stateEmpty   boardState = iota // 0
	stateShip    boardState = iota // 1
	stateHit     boardState = iota // 2
	stateAttempt boardState = iota // 3
)

type battleShipBoard struct {
	Board [10][10]boardState
}

func (b *battleShipBoard) Draw() string {
	str := ""
	str += fmt.Sprint("_|A|B|C|D|E|F|G|H|I|J|_\n")
	for y, stripe := range b.Board {
		str += fmt.Sprintf("%d|", y)
		for _, x := range stripe {
			str += fmt.Sprintf("%s|", x.Draw())
		}
		str += fmt.Sprintf("%d\n", y)
	}
	str += fmt.Sprint("_|A|B|C|D|E|F|G|H|I|J|_\n")
	return str
}

var cblack = ansi.ColorCode("black+h:black")
var cship = ansi.ColorCode("black+h:white")
var chit = ansi.ColorCode("red+h:red")
var cattempt = ansi.ColorCode("yellow:yellow")

func (b boardState) Draw() string {
	if b == stateEmpty {
		return cblack + square + ansi.DefaultBG + ansi.DefaultFG
	}
	if b == stateShip {
		return cship + square + ansi.DefaultBG + ansi.DefaultFG
	}
	if b == stateHit {
		return chit + square + ansi.DefaultBG + ansi.DefaultFG
	}
	if b == stateAttempt {
		return cattempt + square + ansi.DefaultBG + ansi.DefaultFG
	}
	return ""
}

func cordsToNumbers(in string) (X, Y int) {
	in = strings.ToLower(in)

	i1 := int(in[0])
	if i1 > 96 && i1 < 107 {
		X = i1 - 97
	} else {
		return -1, -1
	}

	i2, err := strconv.ParseInt(string(in[1]), 10, 64)
	if err != nil {
		return -1, -1
	}

	if i2 < 11 {
		return X, int(i2)
	}
	return -1, -1
}

func combineBoard(b1, b2 battleShipBoard) string {
	b1r, b2r := b1.Draw(), b2.Draw()

	b1rs, b2rs := strings.Split(b1r, "\n"), strings.Split(b2r, "\n")

	str := ""
	for k, _ := range b1rs {
		str += fmt.Sprintf("%s     %s\n", b1rs[k], b2rs[k])
	}

	return str
}

func makeBoard() battleShipBoard {
	a := battleShipBoard{}
	ri, _ := cr.Int(cr.Reader, big.NewInt(math.MaxInt64))
	rand.Seed(ri.Int64())

	a = placeShip(5, a)
	a = placeShip(4, a)
	a = placeShip(3, a)
	a = placeShip(3, a)
	a = placeShip(2, a)
	return a
}

func placeShip(size int, bo battleShipBoard) battleShipBoard {

	for {
		board := bo
		sideways := rand.Int() % 2

		if sideways == 0 { // ship goes up
			X := rand.Int() % 10
			Y := rand.Int() % 10
			if Y+size > 10 {
				continue
			}

			for y := Y; y < Y+size; y++ {
				if board.Board[y][X] != stateEmpty {
					continue
				}
				board.Board[y][X] = stateShip
			}
			bo = board
			break
		} else {
			X := rand.Int() % 10
			Y := rand.Int() % 10

			if X+size > 10 {
				continue
			}

			for x := X; x < X+size; x++ {
				if board.Board[Y][x] != stateEmpty {
					continue
				}
				board.Board[Y][x] = stateShip
			}
			bo = board
			break
		}
	}

	return bo
}
