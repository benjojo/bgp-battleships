package main

import (
	"bytes"
	"encoding/binary"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"regexp"
	"strconv"
	"strings"

	"github.com/nareix/bits"
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

var templatePath = flag.String("templateFile", "/etc/bird/conf.orig",
	"Where to find the template file")

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

func numberToBitReader(in uint16) *bits.Reader {
	actually := uint16(in)

	bits2 := make([]byte, 2)

	binary.LittleEndian.PutUint16(bits2, actually)
	return &bits.Reader{R: bytes.NewReader(bits2)}
}

var errNotEnoughData = fmt.Errorf("Not enough data to make a move")
var errInvalidType = fmt.Errorf("Invalid community type found")
var errDupeType = fmt.Errorf("Duplicate data read")

func readBGP() (gameIncrementor, X, Y, HitOrMissOnLast int, err error) {
	communities := readCommunities(*monitoredPrefix)

	readCounter, readPosition := false, false

	for _, community := range communities {
		if community.AS == uint16(*communityAS) {
			// okay, so we are now interested!
			r := numberToBitReader(community.Data)
			t, _ := r.ReadBits(2)

			if t == 1 {
				// Counter
				if readCounter {
					// uh we have read it twice, oh dear?
					return 0, 0, 0, 0, errDupeType
				}
				readCounter = true
				c, _ := r.ReadBits(14)
				gameIncrementor = int(c)

			} else if t == 2 {
				if readPosition {
					// uh we have read it twice, oh dear?
					return 0, 0, 0, 0, errDupeType
				}
				readPosition = true
				xp, _ := r.ReadBits(4)
				X = int(xp)
				yp, _ := r.ReadBits(4)
				Y = int(yp)
				hs, _ := r.ReadBits(1)
				HitOrMissOnLast = int(hs)

			} else {
				return 0, 0, 0, 0, errInvalidType
			}
		}
	}

	if readCounter && readPosition {
		return gameIncrementor, X, Y, HitOrMissOnLast, nil
	}
	return 0, 0, 0, 0, errNotEnoughData
}

func writeBGP(gameIncrementor, X, Y, HitOrMissOnLast int) error {

	counternumberbytes := make([]byte, 2)
	counternumberbytesbuffer := bytes.NewBuffer(counternumberbytes)
	counternumberbits := bits.Writer{
		W: counternumberbytesbuffer,
	}

	counternumberbits.WriteBits(1, 2)
	counternumberbits.WriteBits(uint(gameIncrementor), 14)

	counterCommunity := binary.LittleEndian.Uint16(counternumberbytes)

	positionnumberbytes := make([]byte, 2)
	positionnumberbytesbuffer := bytes.NewBuffer(positionnumberbytes)
	positionnumberbits := bits.Writer{
		W: positionnumberbytesbuffer,
	}

	positionnumberbits.WriteBits(2, 2)
	positionnumberbits.WriteBits(uint(X), 4)
	positionnumberbits.WriteBits(uint(Y), 4)
	positionnumberbits.WriteBits(uint(HitOrMissOnLast), 1)

	positionCommunity := binary.LittleEndian.Uint16(positionnumberbytes)

	// Now we have the two community strings counterCommunity and positionCommunity

	templatestring := fmt.Sprintf(
		"\nbgp_community.add((%d,%d));\nbgp_community.add((%d,%d));\n",
		communityAS, positionCommunity, communityAS, counterCommunity)

	templateBytes, err := ioutil.ReadFile(*templatePath)
	if err != nil {
		return err
	}

	birdConfigOutput := strings.Replace(string(templateBytes),
		"###COMMUNITY###", templatestring, 1)

	err = ioutil.WriteFile("/etc/bird/bird.conf", []byte(birdConfigOutput), 0640)
	if err != nil {
		return err
	}

	// now reload bird
	conn, err := net.Dial("unix", "/run/bird/bird.ctl")
	if err != nil {
		log.Fatalf("Unable to connect to bird %s", err.Error())
	}

	defer conn.Close()

	conn.Write([]byte(fmt.Sprintf("reload all\n")))

	buffer := make([]byte, 90000)
	_, err = conn.Read(buffer)

	return err
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
