package state

type Connection struct {
	Username string
	Playername string

	Admin bool
	Remote bool

	PlayerId uint8
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

	Connections map[uint8]Connection
	Players map[uint8]Player
}

func (s *State) AddConnection() {
}

func (s *State) RemoveConnection() {
}

func (s *State) AddPlayer() {
}

func (s *State) RemovePlayer() {
}
