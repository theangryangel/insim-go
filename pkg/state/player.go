package state

type Player struct {
	Playername string
	Plate      string
	Pitting    bool

	ConnectionId uint8

	RacePosition uint8
	RaceLap      uint16
	Speed        uint16 // temporarily disable Mph

	X int32
	Y int32
	Z int32
}
