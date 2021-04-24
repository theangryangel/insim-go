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
		v.Pitting = false
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
		v.Pitting = true
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
