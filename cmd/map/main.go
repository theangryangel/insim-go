// +build ignore

package main

import (
	"flag"
	"fmt"

	"github.com/theangryangel/insim-go/pkg/files"
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

type chat struct {
	SentAt time.Time
	Msg    string

	Plid uint8
	Ucid uint8
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
			data, err := json.Marshal(
				chat{
					SentAt: time.Now(),
					Msg: fmt.Sprintf(
						"%s: %s",
						strings.StripColours(player.Playername),
						strings.StripColours(mso.Msg),
					),
					Plid: mso.Plid,
					Ucid: mso.Ucid,
				},
			)
			if err != nil {
				panic(err)
			}

			s.SendMessage(
				"/events",
				sse.NewMessage(
					"", // id
					string(data),
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
		var data string
		if rst.Racing() {
			data = fmt.Sprintf(
				"Race starting on %s weather=%d,wind=%d,laps=%d",
				rst.Track,
				rst.Weather,
				rst.Wind,
				rst.RaceLaps,
			)
			fmt.Println(data)
		}

		if rst.Qualifying() {
			data = fmt.Sprintf(
				"Qualifying starting on %s weather=%d,wind=%d,duration=%s",
				rst.Track,
				rst.Weather,
				rst.Wind,
				rst.QualifyingDuration(),
			)
			fmt.Println(data)
		}

		chatmsg, err := json.Marshal(
			chat{
				SentAt: time.Now(),
				Msg:    data,
				Plid:   0,
				Ucid:   0,
			},
		)
		if err != nil {
			panic(err)
		}

		s.SendMessage(
			"/events",
			sse.NewMessage(
				"",              // id
				string(chatmsg), // event
				"chat",
			),
		)

		gamedata, err := json.Marshal(client.GameState)
		if err != nil {
			return
		}

		s.SendMessage(
			"/events",
			sse.NewMessage(
				"",               // id
				string(gamedata), // event
				"state",
			),
		)

	})

	c.On(func(client *session.InsimSession, res *protocol.Res) {
		if player, ok := c.GameState.Players[res.Plid]; ok {
			data, err := json.Marshal(
				chat{
					SentAt: time.Now(),
					Msg:    fmt.Sprintf("%s position=%d,btime=%s,ttime=%s", player.Playername, res.ResultNum, res.BestTime(), res.TotalTime()),
				},
			)
			if err != nil {
				panic(err)
			}
			fmt.Println(data)
			s.SendMessage(
				"/events",
				sse.NewMessage(
					"", // id
					string(data),
					"chat", // event
				),
			)
		}
	})

	c.On(func(client *session.InsimSession, con *protocol.Con) {
		a, aok := c.GameState.Players[con.A.Plid]
		b, bok := c.GameState.Players[con.B.Plid]

		if aok && bok {
			data, err := json.Marshal(
				chat{
					SentAt: time.Now(),
					Msg:    fmt.Sprintf("Contact between %s and %s", strings.StripColours(a.Playername), strings.StripColours(b.Playername)),
				},
			)
			if err != nil {
				return
			}
			s.SendMessage(
				"/events",
				sse.NewMessage(
					"", // id
					string(data),
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

	r.Get("/api/track/{code}", func(w http.ResponseWriter, r *http.Request) {
		track := chi.URLParam(r, "code")

		pth := files.Pth{}
		// TODO this is massively unsafe
		pth.Read(fmt.Sprintf("../pth/data/%s.pth", track))

		render.JSON(w, r, pth.FitTo(1024, 1024, 2))
	})

	http.ListenAndServe(":4000", r)
}
