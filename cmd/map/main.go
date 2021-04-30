// +build ignore

package main

import (
	"flag"
	"fmt"

	"github.com/theangryangel/insim-go/pkg/protocol"
	"github.com/theangryangel/insim-go/pkg/session"
	"github.com/theangryangel/insim-go/pkg/strings"

	"encoding/json"
	"time"

	"net/http"

	sse "github.com/alexandrevicenzi/go-sse"

	chi "github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

type playerState struct {
	Plid  uint8
	State interface{}
}

func main() {
	host := flag.String("host", "127.0.0.1:29999", "host:port to dial or hostname if using -relay")
	relay := flag.Bool("relay", false, "Use LFSW relay")
	list := flag.Bool("rlist", false, "Fetch host list from relay")
	flag.Parse()

	dial := *host
	if *relay {
		dial = "isrelay.lfs.net:47474"
	}

	c := session.NewInsimSession()
	defer c.Close()

	// Create SSE server
	s := sse.NewServer(&sse.Options{
		Logger: nil,
	})
	defer s.Shutdown()

	c.On(func(client *session.InsimSession, mci *protocol.Mci) {
		for _, info := range mci.Info {
			if v, ok := c.GameState.Players[info.Plid]; ok {
				data, err := json.Marshal(playerState{Plid: info.Plid, State: v})
				if err != nil {
					continue
				}

				s.SendMessage(
					"/events",
					sse.NewMessage(
						"",             // id
						string(data),   // data
						"player-state", // event
					),
				)
			}
		}
	})

	c.On(func(client *session.InsimSession, pll *protocol.Pll) {
		data, err := json.Marshal(pll)
		if err != nil {
			return
		}

		s.SendMessage(
			"/events",
			sse.NewMessage(
				"", // id
				string(data),
				"player-left", // event
			),
		)
	})

	c.On(func(client *session.InsimSession, mso *protocol.Mso) {
		if player, ok := c.GameState.Players[mso.Plid]; ok {
			data := fmt.Sprintf(
				"%s: %s: %s",
				time.Now(),
				strings.StripColours(player.Playername),
				strings.StripColours(mso.Msg),
			)
			fmt.Println(data)

			s.SendMessage(
				"/events",
				sse.NewMessage(
					"", // id
					data,
					"chat", // event
				),
			)
		}
	})

	c.On(func(client *session.InsimSession, hos *protocol.IrpHos) {
		fmt.Printf("Hosts:\n")
		for _, info := range hos.HInfo {
			fmt.Printf("%s\n", info.HName)
		}
	})

	c.On(func(client *session.InsimSession, err *protocol.IrpErr) {
		fmt.Printf("Relay Error: %s\n", err.ErrMessage())
	})

	c.On(func(client *session.InsimSession, rst *protocol.Rst) {
		if rst.Racing() {
			data := fmt.Sprintf(
				"%s: Race starting on %s weather=%d,wind=%d,laps=%d",
				time.Now(),
				rst.Track,
				rst.Weather,
				rst.Wind,
				rst.RaceLaps,
			)
			fmt.Println(data)
			s.SendMessage(
				"/events",
				sse.NewMessage(
					"",
					data,
					"chat",
				),
			)
		}

		if rst.Qualifying() {
			data := fmt.Sprintf(
				"%s: Qualifying starting on %s weather=%d,wind=%d,duration=%s",
				time.Now(),
				rst.Track,
				rst.Weather,
				rst.Wind,
				rst.QualifyingDuration(),
			)
			fmt.Println(data)
			s.SendMessage(
				"/events",
				sse.NewMessage(
					"",
					data,
					"chat",
				),
			)
		}

		data, err := json.Marshal(client.GameState)
		if err != nil {
			return
		}

		s.SendMessage(
			"/events",
			sse.NewMessage(
				"",           // id
				string(data), // event
				"state",
			),
		)

	})

	c.On(func(client *session.InsimSession, res *protocol.Res) {
		if player, ok := c.GameState.Players[res.Plid]; ok {
			data := fmt.Sprintf("%s: %s position=%d,btime=%s,ttime=%s", time.Now(), player.Playername, res.ResultNum, res.BestTime(), res.TotalTime())
			fmt.Println(data)
			s.SendMessage(
				"/events",
				sse.NewMessage(
					"", // id
					data,
					"chat", // event
				),
			)
		}
	})

	c.On(func(client *session.InsimSession, con *protocol.Con) {
		a, aok := c.GameState.Players[con.A.Plid]
		b, bok := c.GameState.Players[con.B.Plid]

		if aok && bok {
			data := fmt.Sprintf("%s: Contact between %s and %s", time.Now(), strings.StripColours(a.Playername), strings.StripColours(b.Playername))
			s.SendMessage(
				"/events",
				sse.NewMessage(
					"", // id
					data,
					"chat", // event
				),
			)
		}
	})

	go func() {
		fmt.Println("Dialling")
		err := c.Dial(dial)
		if err != nil {
			panic(err)
		}
		defer c.Close()

		fmt.Println("Connected!")

		if *relay && *list {
			hlr := protocol.NewIrpHlr()
			c.Write(hlr)
		} else if *relay {
			c.SelectRelayHost(*host)
			c.RequestState()
		} else {
			c.Init()
			c.RequestState()
		}

		err = c.Scan()

		if err != nil {
			panic(err)
		}
	}()

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.URLFormat)
	r.Use(render.SetContentType(render.ContentTypeJSON))

	fs := http.FileServer(http.Dir("./static"))
	r.Method("GET", "/*", fs)

	r.Mount("/events", s)

	r.Get("/api/connections", func(w http.ResponseWriter, r *http.Request) {
		render.JSON(w, r, c.GameState.Connections)
	})

	r.Get("/api/players", func(w http.ResponseWriter, r *http.Request) {
		render.JSON(w, r, c.GameState.Players)
	})

	r.Get("/api/state", func(w http.ResponseWriter, r *http.Request) {
		render.JSON(w, r, c.GameState)
	})

	http.ListenAndServe(":4000", r)
}
