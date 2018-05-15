package main

import (
	"bytes"
	"encoding/binary"
	"flag"
	"fmt"
	"log"
	"net"
	"regexp"
	"strconv"
	"strings"
)

type bgpCommunity struct {
	AS   uint16
	Data uint16
}

var birdCommunityRegex = regexp.MustCompile(`\((\d+,\d+)\)`)

var monitoredPrefix = flag.String(
	"peerprefix", "1.1.1.0/24", "the prefix of the other side")

var communityAS = flag.Int("communityASN", 23456,
	"The shared community AS used to communicate on")

/*
Three Communities are used:


Type 1: Game counter, Moves up
on each move, so that peers know
when a move happens.

T = Type
Z = Counter number

+-------------------------------+
|T|T|Z|Z|Z|Z|Z|Z|Z|Z|Z|Z|Z|Z|Z|Z|
+-------------------------------+

Type 2: Attack Cords

T = Type
X = X Cords
Y = Y Cords
S = Hit or Miss on last move

+-------------------------------+
|T|T|X|X|X|X|Y|Y|Y|Y|S|-|-|-|-|-|
+-------------------------------+
*/

func numberToByteReader(in uint16) *bytes.Reader {
	actually := uint16(in)

	bits := make([]byte, 2)

	binary.LittleEndian.PutUint16(bits, actually)
	return bytes.NewReader(bits)
}

var errNotEnoughData = fmt.Errorf("Not enough data to make a move")

func readBGP() (gameIncrementor, X, Y, HitOrMissOnLast int, err error) {
	communities := readCommunities(*monitoredPrefix)

	readbits := 0

	for _, community := range communities {
		if community.AS == uint16(*communityAS) {
			// okay, so we are now interested!

		}
	}

	// bits.Reader{}

	return 0, 0, 0, 0
}

func readCommunities(prefix string) (o []bgpCommunity) {
	conn, err := net.Dial("unix", "/run/bird/bird.ctl")
	if err != nil {
		log.Fatalf("Unable to connect to bird %s", err.Error())
	}

	defer conn.Close()

	conn.Write([]byte(fmt.Sprintf("show route all %s\n", prefix)))

	buffer := make([]byte, 90000)
	n, err := conn.Read(buffer)

	if err != nil {
		log.Fatalf("Unable to read from bird %s", err.Error())
	}

	matches :=
		birdCommunityRegex.FindAllStringSubmatch(string(buffer[:n]), -1)

	o = make([]bgpCommunity, 0)

	for _, v := range matches {
		if len(v) == 2 {
			bits := strings.Split(v[1], ",")
			as, _ := strconv.ParseInt(bits[0], 10, 64)
			data, _ := strconv.ParseInt(bits[1], 10, 64)
			o = append(o, bgpCommunity{
				AS:   uint16(as),
				Data: uint16(data),
			})
		}
	}

	// log.Printf("%+v", matches)
	return o
}
