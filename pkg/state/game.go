package state

import (
	"github.com/theangryangel/insim-go/pkg/protocol"
	"sync"
	"time"
)

type GameState struct {
	mu sync.Mutex

	Track              string
	Weather            uint8
	Wind               uint8
	Laps               int32
	Racing             bool
	Voting             bool
	QualifyingDuration time.Duration

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
	if s.Players == nil {
		s.Players = make(map[uint8]*Player)
	}

	id := npl.Plid

	if v, ok := s.Players[id]; ok {
		v.Playername = npl.PName
		v.Plate = npl.Plate
		v.ConnectionId = npl.Ucid
		v.PitGarage = false
	} else {
		s.Players[id] = &Player{
			Playername:   npl.PName,
			Plate:        npl.Plate,
			ConnectionId: npl.Ucid,
		}
	}
	s.mu.Unlock()
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
			v.Speed = info.Speed
			v.RacePosition = info.Position
			v.RaceLap = info.Lap
			v.X = info.X
			v.Y = info.Y
			v.Z = info.Z
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

	s.Track = sta.Track
	s.Weather = sta.Weather
	s.Wind = sta.Wind

	s.Laps = sta.Laps()
	s.Racing = sta.Racing()
	s.QualifyingDuration = sta.QualifyingDuration()
}

func (s *GameState) FromRst(rst *protocol.Rst) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.Track = rst.Track
	s.Weather = rst.Weather
	s.Wind = rst.Wind

	s.Laps = rst.Laps()
	s.Racing = rst.Racing()
	s.QualifyingDuration = rst.QualifyingDuration()
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
