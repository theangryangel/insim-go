package state

import (
	"time"

	"github.com/theangryangel/insim-go/pkg/facts"
)

type Position struct {
	Speed uint16
	X     int32
	Y     int32
	Z     int32
}

type Gaps struct {
	// list of most recent splits 0=1, 1=2, 2=3, 4=start/finish
	// used to help populate GTime
	Duration [facts.MaxSplitCount + 1]time.Duration `json:"-"`

	// lap number by split index.
	// i.e. last SPX, split = 0, lap was N, split = 1, lap was N-1
	Lap [facts.MaxSplitCount + 1]uint16 `json:"-"`

	Leader string // TODO

	Next string // TODO it would be nice if this could be not just a string
	Prev string
}

func (g *Gaps) Update(spx uint8, lap uint16, etime time.Duration) {
	g.Duration[spx] = etime
	g.Lap[spx] = lap
}

func (g *Gaps) Reset() {
	for i := 0; i < facts.MaxSplitCount+1; i++ {
		g.Duration[i] = time.Duration(0)
		g.Lap[i] = 0
	}

	g.Next = ""
	g.Prev = ""
	g.Leader = ""
}

type Player struct {
	Playername string
	Plate      string

	PitGarage bool
	PitLane   bool

	ConnectionId uint8

	RacePosition uint8
	RaceLap      uint16
	RaceFinished bool

	NumStops uint32

	Vehicle string

	TTime time.Duration // TODO json encode a string
	BTime time.Duration

	Gaps     Gaps
	Position Position

	YellowFlag bool
	BlueFlag   bool
}

func (p *Player) Reset() {
	p.PitGarage = false
	p.PitLane = false
	p.RaceFinished = false
	p.TTime = time.Duration(0)
	p.BTime = time.Duration(0)
	p.NumStops = 0
	p.Gaps.Reset()
}
