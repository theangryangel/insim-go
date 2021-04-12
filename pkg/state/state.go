package state

type Connection struct {
	Username string
	Playername string

	Admin bool
	Remote bool
}

type Player struct {
	Playername string
	Plate string

	ConnectionId uint8
}

type State struct {
	Track string
	Weather uint8
	Wind uint8

	Connections map[uint8]*Connection
	Players map[uint8]*Player
}

func (s *State) UpdateConnection(id uint8, username string, playername string, admin bool, remote bool) {
	// TODO: Needs locks
	if s.Connections == nil {
		s.Connections = make(map[uint8]*Connection)
	}

	if v, ok := s.Connections[id]; ok {
		v.Username = username
		v.Playername = playername
		v.Admin = admin
		v.Remote = remote
	} else{
		s.Connections[id] = &Connection{
			Username: username,
			Playername: playername,
			Admin: admin,
			Remote: remote,
		}
	}
}

func (s *State) RemoveConnection(id uint8) {
	if s.Connections == nil {
		return
	}

	if _, ok := s.Connections[id]; ok {
		delete(s.Connections, id)
	}
}

func (s *State) UpdatePlayer(id uint8, playername string, ucid uint8) {
	// TODO: Needs locks
	if s.Players == nil {
		s.Players = make(map[uint8]*Player)
	}

	if v, ok := s.Players[id]; ok {
		v.Playername = playername
		v.ConnectionId = ucid
	} else{
		s.Players[id] = &Player{
			Playername: playername,
			ConnectionId: ucid,
		}
	}
}

func (s *State) RemovePlayer(id uint8) {
	if s.Players == nil {
		return
	}

	if _, ok := s.Players[id]; ok {
		delete(s.Players, id)
	}
}
