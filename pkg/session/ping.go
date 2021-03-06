package session

import (
	"fmt"
	"github.com/theangryangel/insim-go/pkg/protocol"
)

func usePing(c *InsimSession) {
	c.On(func(client *InsimSession, data *protocol.Tiny) {
		if data.SubT == 0 {
			client.Write(data)
			fmt.Println("Ping? Pong!")
		}
	})
}
