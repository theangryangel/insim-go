package state

type Player struct {
	Playername string
	Plate      string
	PitGarage  bool
	PitLane    bool

	ConnectionId uint8

	RacePosition uint8
	RaceLap      uint16
	Speed        uint16 // temporarily disable Mph

	YellowFlag bool
	BlueFlag   bool

	X int32
	Y int32
	Z int32
}
