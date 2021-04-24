package session

import (
	"github.com/theangryangel/insim-go/pkg/protocol"
)

func useGameState(c *InsimSession) {
	c.On(func(c *InsimSession, ncn *protocol.Ncn) {
		c.GameState.FromNcn(ncn)
	})

	c.On(func(c *InsimSession, cnl *protocol.Cnl) {
		c.GameState.FromCnl(cnl)
	})

	c.On(func(c *InsimSession, npl *protocol.Npl) {
		c.GameState.FromNpl(npl)
	})

	c.On(func(c *InsimSession, pll *protocol.Pll) {
		c.GameState.FromPll(pll)
	})

	c.On(func(c *InsimSession, pll *protocol.Plp) {
		c.GameState.FromPlp(pll)
	})

	c.On(func(c *InsimSession, mci *protocol.Mci) {
		c.GameState.FromMci(mci)
	})

	c.On(func(c *InsimSession, sta *protocol.Sta) {
		c.GameState.Track = sta.Track
		c.GameState.Weather = sta.Weather
		c.GameState.Wind = sta.Wind

		c.GameState.Laps = sta.Laps()
		c.GameState.Racing = sta.Racing()
		c.GameState.QualifyingDuration = sta.QualifyingDuration()
	})
}
