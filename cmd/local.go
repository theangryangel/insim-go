// +build ignore

package main

import (
	"github.com/theangryangel/insim-go/pkg/session"
	"github.com/theangryangel/insim-go/pkg/protocol"
	"fmt"
	"flag"
)

func main() {
	host := flag.String("host", "127.0.0.1:29999", "host:port to dial or hostname if using -relay")
	relay := flag.Bool("relay", false, "Use LFSW relay")
	flag.Parse()

	dial := *host
	if *relay {
		dial = "isrelay.lfs.net:47474"
	}

	c := session.NewInsimSession()

	c.On(func(client *session.InsimSession, ver *protocol.Ver) {
		if ver.ReqI != 1 {
			return
		}
		fmt.Printf("Connected. Ver = %v\n", ver)

	})

	c.On(func(client *session.InsimSession, mso *protocol.Mso) {
		fmt.Printf("Msg: %s\n", mso.Msg)
	})

	c.On(func(client *session.InsimSession, data *protocol.Tiny) {
		if data.SubT == 0 {
			client.Write(data)
			fmt.Println("Ping? Pong!")
		}
	})

	c.On(func(client *session.InsimSession, sta *protocol.Sta) {
		fmt.Printf("Track: %s\n", sta.Track)
	})

	err := c.Dial(dial)
	if err != nil {
		panic(err)
	}
	defer c.Close()

	fmt.Println("Connected!")

	if *relay {
		c.SelectRelayHost(*host)
	} else {
		c.Init()
	}

	c.RequestState()

	for {
	  err := c.Read()
		if err != nil {

			if err == session.ErrUnknownType {
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
