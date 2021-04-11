package session

import (
	"fmt"
	"github.com/theangryangel/insim-go/pkg/protocol"
)

func (c *InsimSession) trackGameState() {
	c.On(func (c *InsimSession, sta *protocol.Sta) {
		c.gamestate.Track = sta.Track
		c.gamestate.Weather = sta.Weather
		fmt.Println("Changed game state")
	})
}
