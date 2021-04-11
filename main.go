// +build ignore

package main

import (
	"github.com/theangryangel/insim/pkg/session"
	"github.com/theangryangel/insim/pkg/protocol"
	"fmt"
)

func main() {
	c := session.NewInsimSession()
	err := c.Dial("192.168.0.250:29999")
	if err != nil {
		panic(err)
	}
	defer c.Close()

	fmt.Println("Connected!")

	c.On(func(client *session.InsimSession, mso *protocol.Mso) {
		fmt.Printf("Handled Msg: %s\n", mso.Msg)
	})

	c.On(func(client *session.InsimSession, data *protocol.Tiny) {
		if data.SubT == 0 {
			client.Write(data)
			fmt.Println("Ping? Pong!")
		}
	})

	c.On(func(client *session.InsimSession, sta *protocol.Sta) {
		fmt.Printf("STA: Track: %s\n", sta, sta.Track)
	})

	for {
	  err := c.Read()
		if err != nil {

			if err == session.ErrUnknownType {
				fmt.Println("Unknown Packet")
				continue
			}

			if err == session.ErrNotEnough {
				fmt.Println("Not enough data")
				continue
			}

			panic(err)
		}
	}
}
