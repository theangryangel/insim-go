package state

import (
	"fmt"
	"github.com/theangryangel/insim-go/pkg/facts"
	"github.com/theangryangel/insim-go/pkg/protocol"
	"sync"
	"time"
)

type GameState struct {
	mu sync.Mutex `json:"-"`

	// TODO are these more like something an event state?
	// probably. What does that look like?

	// Fastest time/player?
	Track              *facts.Track
	Weather            uint8
	Wind               uint8
	Laps               int32
	Racing             bool
	Voting             bool
	QualifyingDuration time.Duration

	// TODO could do with a ConnectionList and PlayerList type to move some of these functions info
	// Would allow less locking of the whole GameState
	Connections map[uint8]*Connection
	Players     map[uint8]*Player
}

func (s *GameState) FromNcn(ncn *protocol.Ncn) {
	s.mu.Lock()

	if s.Connections == nil {
		s.Connections = make(map[uint8]*Connection)
	}

	id := ncn.Ucid

	if v, ok := s.Connections[id]; ok {
		v.Username = ncn.UName
		v.Playername = ncn.PName
		v.Admin = ncn.IsAdmin()
		v.Remote = ncn.IsRemote()
	} else {
		s.Connections[id] = &Connection{
			Username:   ncn.UName,
			Playername: ncn.PName,
			Admin:      ncn.IsAdmin(),
			Remote:     ncn.IsRemote(),
		}
	}

	s.mu.Unlock()
}

func (s *GameState) FromCnl(cnl *protocol.Cnl) {
	s.mu.Lock()
	if s.Connections == nil {
		return
	}

	id := cnl.Ucid

	if _, ok := s.Connections[id]; ok {
		delete(s.Connections, id)
	}
	s.mu.Unlock()
}

func (s *GameState) FromNpl(
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

func (s *GameState) FromPll(pll *protocol.Pll) {
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

func (s *GameState) FromPlp(plp *protocol.Plp) {
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

func (s *GameState) FromMci(mci *protocol.Mci) {
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

func (s *GameState) FromToc(toc *protocol.Toc) {
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

func (s *GameState) FromFlg(flg *protocol.Flg) {
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

func (s *GameState) FromSta(sta *protocol.Sta) {
	s.mu.Lock()
	defer s.mu.Unlock()

	track, err := facts.TrackFromCode(sta.Track)
	if err != nil {
		panic(err)
	}

	s.Track = track
	s.Weather = sta.Weather
	s.Wind = sta.Wind

	s.Laps = sta.Laps()
	s.Racing = sta.Racing()
	s.QualifyingDuration = sta.QualifyingDuration()
}

func (s *GameState) FromRst(rst *protocol.Rst) {
	s.mu.Lock()
	defer s.mu.Unlock()

	track, err := facts.TrackFromCode(rst.Track)
	if err != nil {
		panic(err)
	}

	s.Track = track
	s.Weather = rst.Weather
	s.Wind = rst.Wind

	s.Laps = rst.Laps()
	s.Racing = rst.Racing()
	s.QualifyingDuration = rst.QualifyingDuration()

	for _, p := range s.Players {
		p.Reset()
	}
}

func (s *GameState) FromPla(pla *protocol.Pla) {
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

func (s *GameState) FromVtn(vtn *protocol.Vtn) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.Voting = vtn.Voting()
}

func (s *GameState) FromTiny(tiny *protocol.Tiny) {

	s.mu.Lock()
	defer s.mu.Unlock()

	if tiny.SubT == protocol.TINY_VTC {
		s.Voting = false
	}
}

func (s *GameState) FromSmall(small *protocol.Small) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if small.SubT == protocol.SMALL_VTA {
		s.Voting = false
	}
}

func (s *GameState) FromFin(fin *protocol.Fin) {
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

func (s *GameState) FromRes(fin *protocol.Res) {
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

func (s *GameState) FromLap(lap *protocol.Lap) {
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

func (s *GameState) FromSpx(spx *protocol.Spx) {
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

func (s *GameState) UpdateGaps(plid uint8, spx uint8) {
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
