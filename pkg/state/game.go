package state

import (
	"encoding/json"
	"sync"
	"time"

	"github.com/theangryangel/insim-go/pkg/facts"
	"github.com/theangryangel/insim-go/pkg/protocol"
)

type Event struct {
	mu sync.Mutex

	// Fastest time/player?
	Track              *facts.Track
	Weather            uint8
	Wind               uint8
	Laps               int32
	Racing             bool
	Voting             bool
	QualifyingDuration time.Duration

	MultiPlayer bool
	HostName    string
}

type GameState struct {
	Event Event

	Connections ConnectionList
	Players     PlayerList
}

func (s *GameState) UnmarshalJSON(b []byte) error {
	return nil
}

func (s *GameState) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Event       *Event
		Players     map[uint8]*Player
		Connections map[uint8]*Connection
	}{
		Event:       &s.Event,
		Players:     s.Players.Players,
		Connections: s.Connections.Connections,
	})
}

func (s *GameState) FromSta(sta *protocol.Sta) {
	s.Event.FromSta(sta)
}

func (s *Event) FromSta(sta *protocol.Sta) {
	s.mu.Lock()
	defer s.mu.Unlock()

	track, err := facts.TrackFromCode(sta.Track[0:3])
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
	s.Event.FromRst(rst)
	s.Players.Reset()
}

func (s *Event) FromRst(rst *protocol.Rst) {
	s.mu.Lock()
	defer s.mu.Unlock()

	track, err := facts.TrackFromCode(rst.Track[0:3])
	if err != nil {
		panic(err)
	}

	s.Track = track
	s.Weather = rst.Weather
	s.Wind = rst.Wind

	s.Laps = rst.Laps()
	s.Racing = rst.Racing()
	s.QualifyingDuration = rst.QualifyingDuration()
}

func (s *GameState) FromVtn(vtn *protocol.Vtn) {
	s.Event.FromVtn(vtn)
}

func (s *Event) FromVtn(vtn *protocol.Vtn) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.Voting = vtn.Voting()
}

func (s *GameState) FromTiny(tiny *protocol.Tiny) {
	s.Event.mu.Lock()
	defer s.Event.mu.Unlock()

	if tiny.SubT == protocol.TINY_VTC {
		s.Event.Voting = false
	}
}

func (s *GameState) FromSmall(small *protocol.Small) {
	s.Event.mu.Lock()
	defer s.Event.mu.Unlock()

	if small.SubT == protocol.SMALL_VTA {
		s.Event.Voting = false
	}
}

func (s *GameState) FromNcn(ncn *protocol.Ncn) {
	s.Connections.FromNcn(ncn)
}

func (s *GameState) FromCnl(cnl *protocol.Cnl) {
	s.Connections.FromCnl(cnl)
}

func (s *GameState) FromNpl(
	npl *protocol.Npl,
) {
	s.Players.FromNpl(npl)
}

func (s *GameState) FromPll(pll *protocol.Pll) {
	s.Players.FromPll(pll)
}

func (s *GameState) FromPlp(plp *protocol.Plp) {
	s.Players.FromPlp(plp)
}

func (s *GameState) FromMci(mci *protocol.Mci) {
	s.Players.FromMci(mci)
}

func (s *GameState) FromToc(toc *protocol.Toc) {
	s.Players.FromToc(toc)
}

func (s *GameState) FromFlg(flg *protocol.Flg) {
	s.Players.FromFlg(flg)
}

func (s *GameState) FromPla(pla *protocol.Pla) {
	s.Players.FromPla(pla)
}

func (s *GameState) FromFin(fin *protocol.Fin) {
	s.Players.FromFin(fin)
}

func (s *GameState) FromRes(res *protocol.Res) {
	s.Players.FromRes(res)
}

func (s *GameState) FromLap(lap *protocol.Lap) {
	s.Players.FromLap(lap)
}

func (s *GameState) FromSpx(spx *protocol.Spx) {
	s.Players.FromSpx(spx)
}

func (s *GameState) UpdateGaps(plid uint8, spx uint8) {
	s.Players.UpdateGaps(plid, spx)
}

func (s *GameState) FromIsm(ism *protocol.Ism) {
	s.Event.mu.Lock()
	defer s.Event.mu.Unlock()
	s.Event.MultiPlayer = true
	s.Event.HostName = ism.HName
}
