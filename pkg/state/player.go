package state

import (
	"fmt"
	"sync"
	"time"

	"github.com/theangryangel/insim-go/pkg/facts"
	"github.com/theangryangel/insim-go/pkg/protocol"
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
		g.Lap[i] = 1
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

func (s *PlayerList) Get(plid uint8) (*Player, bool) {
	player, ok := s.Players[plid]
	return player, ok
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

type PlayerList struct {
	mu sync.Mutex

	Players map[uint8]*Player
}

func (s *PlayerList) Reset() {
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, p := range s.Players {
		p.Reset()
	}
}

func (s *PlayerList) FromNpl(
	npl *protocol.Npl,
) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.Players == nil {
		s.Players = make(map[uint8]*Player)
	}

	id := npl.Plid

	if v, ok := s.Players[id]; ok {
		v.Playername = npl.PName
		v.Plate = npl.Plate
		v.ConnectionId = npl.Ucid
		v.Vehicle = npl.CName
		v.Reset()
	} else {
		s.Players[id] = &Player{
			Playername:   npl.PName,
			Plate:        npl.Plate,
			ConnectionId: npl.Ucid,
			Vehicle:      npl.CName,
		}
	}
}

func (s *PlayerList) FromPll(pll *protocol.Pll) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.Players == nil {
		return
	}

	id := pll.Plid

	if _, ok := s.Players[id]; ok {
		delete(s.Players, id)
	}
}

func (s *PlayerList) FromPlp(plp *protocol.Plp) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.Players == nil {
		return
	}

	id := plp.Plid

	if v, ok := s.Players[id]; ok {
		v.PitGarage = true
		v.PitLane = false
		if !v.RaceFinished {
			v.RacePosition = 0
		}
	}
}

func (s *PlayerList) FromMci(mci *protocol.Mci) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.Players == nil {
		return
	}

	for _, info := range mci.Info {

		if v, ok := s.Players[info.Plid]; ok {
			if !v.RaceFinished {
				v.RacePosition = info.Position
				v.RaceLap = info.Lap
			}

			v.Position.Speed = info.Speed
			v.Position.X = info.X
			v.Position.Y = info.Y
			v.Position.Z = info.Z
		}
	}
}

func (s *PlayerList) FromToc(toc *protocol.Toc) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.Players == nil {
		return
	}

	id := toc.Plid

	if v, ok := s.Players[id]; ok {
		v.ConnectionId = toc.NewUcid
	}
}

func (s *PlayerList) FromFlg(flg *protocol.Flg) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.Players == nil {
		return
	}

	id := flg.Plid

	if v, ok := s.Players[id]; ok {
		flag, on := flg.Changed()

		v.YellowFlag = (flag == protocol.FLG_YELLOW && on)
		v.BlueFlag = (flag == protocol.FLG_BLUE && on)
	}
}

func (s *PlayerList) FromPla(pla *protocol.Pla) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.Players == nil {
		return
	}

	id := pla.Plid

	if v, ok := s.Players[id]; ok {
		v.PitLane = pla.Entering()
	}
}

func (s *PlayerList) FromFin(fin *protocol.Fin) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.Players == nil {
		return
	}

	id := fin.Plid

	if v, ok := s.Players[id]; ok {
		v.RaceFinished = true
		v.TTime = fin.TotalTime()
		v.BTime = fin.BestTime()
	}
}

func (s *PlayerList) FromRes(fin *protocol.Res) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.Players == nil {
		return
	}

	id := fin.Plid

	if v, ok := s.Players[id]; ok {
		v.RaceFinished = true
		v.RacePosition = fin.RacePosition()
		v.TTime = fin.TotalTime()
		v.BTime = fin.BestTime()
	}
}

func (s *PlayerList) FromLap(lap *protocol.Lap) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.Players == nil {
		return
	}

	id := lap.Plid

	if v, ok := s.Players[id]; ok {
		if !v.RaceFinished {
			v.TTime = lap.ElapsedTime()

			if v.BTime.Nanoseconds() <= 0 || lap.LapTime() < v.BTime {
				v.BTime = lap.LapTime()
			}

			// we consider the lap to be the 4th "split"
			v.Gaps.Update(facts.MaxSplitCount, v.RaceLap, lap.ElapsedTime())
			s.UpdateGaps(id, facts.MaxSplitCount)

		}
		v.NumStops = uint32(lap.NumStops)
	}
}

func (s *PlayerList) FromSpx(spx *protocol.Spx) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.Players == nil {
		return
	}

	id := spx.Plid

	if v, ok := s.Players[id]; ok {
		v.Gaps.Update(spx.Split-1, v.RaceLap, spx.ElapsedTime())
		s.UpdateGaps(id, spx.Split-1)
	}
}

func (s *PlayerList) UpdateGaps(plid uint8, spx uint8) {
	// TODO Fix "65534 laps" bug
	// probably just a logic error around the "4th" split or "lap"

	player, ok := s.Players[plid]
	if !ok {
		return
	}

	if player.RacePosition == 0 || player.RaceFinished {
		// we don't update the gap if they're fresh out of the garage or they're finished the race
		return
	}

	for oplid, oplayer := range s.Players {
		if oplid == plid {
			continue
		}

		if player.RacePosition == oplayer.RacePosition+1 {
			// find the player in front of us

			if player.Gaps.Lap[spx] == oplayer.Gaps.Lap[spx] {
				// are they are on the same lap?
				gap := (oplayer.Gaps.Duration[spx] - player.Gaps.Duration[spx])
				player.Gaps.Next = gap.String()
				oplayer.Gaps.Prev = (-1 * gap).String()
			} else {
				gap := int32(player.Gaps.Lap[spx] - oplayer.Gaps.Lap[spx])
				player.Gaps.Next = fmt.Sprintf("%d laps", gap)
				oplayer.Gaps.Prev = fmt.Sprintf("%d laps", -1*gap)
			}

			return
		}
	}

	player.Gaps.Next = ""
	player.Gaps.Prev = ""
}
