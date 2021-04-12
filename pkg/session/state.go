package session

import (
	"fmt"
	"github.com/theangryangel/insim-go/pkg/protocol"
)

func (c *InsimSession) trackGameState() {
	c.On(func (c *InsimSession, ncn *protocol.Ncn) {
		c.GameState.UpdateConnection(
			ncn.Ucid,
			ncn.UName,
			ncn.PName,
			ncn.Admin == 1, // TODO: Add IsAdmin() support to ncn
			true, // TODO: Add flags support to Ncn
		)

		for ucid, conn := range c.GameState.Connections {
			fmt.Printf("Conn: %d = %v\n", ucid, conn)
		}
	})

	c.On(func (c *InsimSession, cnl *protocol.Cnl) {
		c.GameState.RemoveConnection(cnl.Ucid)
		fmt.Printf("Connection Left %d\n", cnl.Ucid)
	})

	c.On(func (c *InsimSession, npl *protocol.Npl) {
		c.GameState.UpdatePlayer(npl.Plid, npl.PName, npl.Ucid)
		for plid, ply := range c.GameState.Players {
			fmt.Printf("Player: %d = %v\n", plid, ply)
		}
	})

	c.On(func (c *InsimSession, pll *protocol.Pll) {
		c.GameState.RemovePlayer(pll.Plid)
		fmt.Printf("Player %d left\n", pll.Plid)
	})

	c.On(func (c *InsimSession, sta *protocol.Sta) {
		c.GameState.Track = sta.Track
		c.GameState.Weather = sta.Weather
		fmt.Println("Changed game state")
	})
}
