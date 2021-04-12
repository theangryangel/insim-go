package state

import (
	"time"
	"github.com/theangryangel/insim-go/pkg/protocol"
)

type Connection struct {
	Username string
	Playername string

	Admin bool
	Remote bool
}

type Player struct {
	Playername string
	Plate string
	Pitting bool

	ConnectionId uint8

	RacePosition uint8
	RaceLap uint16
	SpeedMph float32
}

type State struct {
	Track string
	Weather uint8
	Wind uint8
	Laps uint8
	Racing bool
	QualifyingDuration time.Duration

	Connections map[uint8]*Connection
	Players map[uint8]*Player
}

func (s *State) FromNcn(ncn *protocol.Ncn) {
	// TODO: Needs locks
	if s.Connections == nil {
		s.Connections = make(map[uint8]*Connection)
	}

	id := ncn.Ucid

	if v, ok := s.Connections[id]; ok {
		v.Username = ncn.UName
		v.Playername = ncn.PName
		v.Admin = ncn.IsAdmin()
		v.Remote = ncn.IsRemote()
	} else{
		s.Connections[id] = &Connection{
			Username: ncn.UName,
			Playername: ncn.PName,
			Admin: ncn.IsAdmin(),
			Remote: ncn.IsRemote(),
		}
	}
}

func (s *State) FromCnl(cnl *protocol.Cnl) {
	if s.Connections == nil {
		return
	}

	id := cnl.Ucid

	if _, ok := s.Connections[id]; ok {
		delete(s.Connections, id)
	}
}

func (s *State) FromNpl(
	npl *protocol.Npl,
) {
	// TODO: Needs locks
	if s.Players == nil {
		s.Players = make(map[uint8]*Player)
	}

	id := npl.Plid

	if v, ok := s.Players[id]; ok {
		v.Playername = npl.PName
		v.Plate = npl.Plate
		v.ConnectionId = npl.Ucid
		v.Pitting = false
	} else{
		s.Players[id] = &Player{
			Playername: npl.PName,
			Plate: npl.Plate,
			ConnectionId: npl.Ucid,
		}
	}
}

func (s *State) FromPll(pll *protocol.Pll){
	if s.Players == nil {
		return
	}

	id := pll.Plid

	if _, ok := s.Players[id]; ok {
		delete(s.Players, id)
	}
}

func (s *State) FromPlp(plp *protocol.Plp){
	if s.Players == nil {
		return
	}

	id := plp.Plid

	if v, ok := s.Players[id]; ok {
		v.Pitting = true
	}
}

func (s *State) FromMci(mci *protocol.Mci) {
	if s.Players == nil {
		return
	}

	for _, info := range mci.Info {

		if v, ok := s.Players[info.Plid]; ok {
			v.SpeedMph = info.Mph()
			v.RacePosition = info.Position
			v.RaceLap = info.Lap
		}
	}

}
