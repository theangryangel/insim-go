package state

import (
	"fmt"
	"sync"
	"time"

	"github.com/theangryangel/insim-go/pkg/facts"
	"github.com/theangryangel/insim-go/pkg/geometry"
	"github.com/theangryangel/insim-go/pkg/protocol"
)

type Split struct {
	Time  time.Duration // time for this lap
	ETime time.Duration // total duration

	RacePosition uint8
}

type Lap struct {
	Split [facts.MaxSplitCount + 1]*Split

	Time  time.Duration
	ETime time.Duration
}

type Player struct {
	Playername string
	Plate      string

	PitGarage bool
	PitLane   bool

	ConnectionId uint8

	RaceStartPosition uint8
	RacePosition      uint8
	RaceLap           uint16
	RaceFinished      bool

	NumStops uint32

	Vehicle string

	TTime time.Duration
	LTime time.Duration
	BTime time.Duration

	// record a map of the splits
	LapTimings map[uint16]*Lap

	CurrentLapTiming *Lap

	// record the lap number the last time we got a split
	// i.e. split 1 on lap 1 is recorded as LastTimingLap[0] = 1
	lapTimingSpx [facts.MaxSplitCount + 1]uint16

	GapNext string
	GapPrev string

	Position geometry.FixedPoint
	Speed    uint16

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
	p.LTime = time.Duration(0)
	p.NumStops = 0
	p.LapTimings = make(map[uint16]*Lap, 0)
	p.CurrentLapTiming = nil
	p.GapNext = ""
	p.GapPrev = ""
	for i := 0; i < facts.MaxSplitCount; i++ {
		p.lapTimingSpx[i] = 1
	}
	p.RaceLap = 1
}

func (p *Player) GetLapTiming(lap uint16) *Lap {
	l, ok := p.LapTimings[lap]
	if !ok {
		l = &Lap{}
		p.LapTimings[lap] = l
	}

	return l
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
		s.Players[id].Reset()
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
			v.Reset()
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

			v.Speed = info.Speed
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
		v.RaceLap = fin.LapsDone
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
		v.RaceLap = fin.LapsDone
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
		l := v.GetLapTiming(lap.LapsDone)
		l.Time = lap.LapTime()

		// record the lap as our "last split"
		// so we can use this for our gaps calc
		l.Split[facts.MaxSplitCount] = &Split{
			Time:         lap.LapTime(),
			ETime:        lap.ElapsedTime(),
			RacePosition: v.RacePosition,
		}

		v.CurrentLapTiming = l
		v.lapTimingSpx[facts.MaxSplitCount] = lap.LapsDone

		if !v.RaceFinished {
			v.LTime = lap.LapTime()
			v.TTime = lap.ElapsedTime()

			if v.BTime.Nanoseconds() <= 0 || lap.LapTime() < v.BTime {
				v.BTime = lap.LapTime()
			}

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
		l := v.GetLapTiming(v.RaceLap)
		v.CurrentLapTiming = l

		l.Split[spx.Split-1] = &Split{
			Time:         spx.SplitTime(),
			ETime:        spx.ElapsedTime(),
			RacePosition: v.RacePosition,
		}

		v.lapTimingSpx[spx.Split-1] = v.RaceLap

		s.UpdateGaps(id, spx.Split-1)
	}
}

func (s *PlayerList) UpdateGaps(plid uint8, spx uint8) {
	player, ok := s.Get(plid)
	if !ok {
		return
	}

	if player.RacePosition == 1 {
		player.GapNext = ""
	}

	if player.RacePosition == 0 || player.RaceFinished {
		// we don't update the gap if they're fresh out of the garage or they're finished the race
		return
	}

	if player.LapTimings[player.RaceLap] == nil || player.LapTimings[player.RaceLap].Split[spx] == nil {
		// we're missing data. probably joined half way through a race.
		// bail out early.
		return
	}

	for oplid, oplayer := range s.Players {
		if oplid == plid {
			continue
		}

		// find the player in front of us
		if player.RacePosition == oplayer.RacePosition+1 {
			if oplayer.LapTimings[player.RaceLap] == nil || oplayer.LapTimings[player.RaceLap].Split[spx] == nil {
				// missing data, probably joined midway through a race.
				return
			}

			// are they are on the same lap?
			if player.lapTimingSpx[spx] == oplayer.lapTimingSpx[spx] {
				gap := (oplayer.LapTimings[player.RaceLap].Split[spx].ETime - player.LapTimings[player.RaceLap].Split[spx].ETime)
				player.GapNext = gap.String()
				oplayer.GapPrev = (-1 * gap).String()
			} else {
				gap := int32(player.RaceLap) - int32(oplayer.RaceLap)
				player.GapNext = fmt.Sprintf("%d laps", gap)
				oplayer.GapPrev = fmt.Sprintf("%d laps", -1*gap)
			}

			return
		}
	}

	player.GapNext = ""
	player.GapPrev = ""
}

func (s *PlayerList) FromReo(reo *protocol.Reo) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.Players == nil {
		return
	}

	fmt.Printf("REO: nump=%d players=%v\n", reo.NumP, reo.Plid)

	for idx := 0; idx < int(reo.NumP); idx++ {
		if p, ok := s.Get(reo.Plid[idx]); ok {
			p.RaceStartPosition = uint8(idx) + 1
		}
	}
}
