package state

import (
	"time"
)

type Player struct {
	Playername string
	Plate      string

	PitGarage bool
	PitLane   bool

	ConnectionId uint8

	RacePosition uint8
	RaceLap      uint16
	RaceFinished bool
	Speed        uint16 // temporarily disable Mph

	NumStops uint32

	Vehicle string

	TTime time.Duration // TODO json encode a string
	BTime time.Duration

	YellowFlag bool
	BlueFlag   bool

	X int32
	Y int32
	Z int32
}
