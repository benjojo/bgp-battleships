package main

import (
	"fmt"
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

func (b *battleShipBoard) Draw() {
	fmt.Print("_|A|B|C|D|E|F|G|H|I|J|_\n")
	for y, stripe := range b.Board {
		fmt.Printf("%d|", y)
		for _, x := range stripe {
			fmt.Printf("%s|", x.Draw())
		}
		fmt.Printf("%d\n", y)
	}
	fmt.Print("_|A|B|C|D|E|F|G|H|I|J|_\n")
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
