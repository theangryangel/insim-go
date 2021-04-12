// +build ignore

package main

import (
	"flag"
	"fmt"

	"github.com/theangryangel/insim-go/pkg/protocol"
	"github.com/theangryangel/insim-go/pkg/session"

	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
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

	c.On(func(client *session.InsimSession, mso *protocol.Mso) {
		fmt.Printf("Msg: %s\n", mso.Msg)
	})

	go func() {
		fmt.Println("Dialling")
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
		err = c.ReadLoop()
		if err != nil {
			panic(err)
		}
	}()

  // Echo instance
  e := echo.New()

  // Middleware
  e.Use(middleware.Logger())
  e.Use(middleware.Recover())

  e.GET("/api/connections", func(ec echo.Context) error {
		return ec.JSON(http.StatusOK, c.GameState.Connections)
	})

  e.GET("/api/players", func(ec echo.Context) error {
		return ec.JSON(http.StatusOK, c.GameState.Players)
	})

  e.GET("/api/state", func(ec echo.Context) error {
		return ec.JSON(http.StatusOK, c.GameState)
	})

	e.Static("/", "static")

  e.Logger.Fatal(e.Start("0.0.0.0:4000"))
}
