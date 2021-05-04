package state

import (
	"github.com/theangryangel/insim-go/pkg/protocol"
	"sync"
)

type Connection struct {
	Username   string
	Playername string

	Admin  bool
	Remote bool
}

type ConnectionList struct {
	mu sync.Mutex `json:"-"`

	Connections map[uint8]*Connection
}

func (s *ConnectionList) Get(ucid uint8) (*Connection, bool) {
	connection, ok := s.Connections[ucid]
	return connection, ok
}

func (s *ConnectionList) FromNcn(ncn *protocol.Ncn) {
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
}

func (s *ConnectionList) FromCnl(cnl *protocol.Cnl) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.Connections == nil {
		return
	}

	id := cnl.Ucid

	if _, ok := s.Connections[id]; ok {
		delete(s.Connections, id)
	}
}
