// +build ignore

package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/theangryangel/insim-go/pkg/files"
	"github.com/theangryangel/insim-go/pkg/protocol"
	"github.com/theangryangel/insim-go/pkg/session"
	"github.com/theangryangel/insim-go/pkg/strings"

	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"path/filepath"
	"time"

	sse "github.com/alexandrevicenzi/go-sse"

	chi "github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"

	svg "github.com/ajstarks/svgo/float"
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

func CleanPath(path string) string {
	// Deal with empty strings nicely.
	if path == "" {
		return ""
	}

	// Ensure that all paths are cleaned (especially problematic ones like
	// "/../../../../../" which can cause lots of issues).
	path = filepath.Clean(path)

	// If the path isn't absolute, we need to do more processing to fix paths
	// such as "../../../../<etc>/some/path". We also shouldn't convert absolute
	// paths to relative ones.
	if !filepath.IsAbs(path) {
		path = filepath.Clean(string(os.PathSeparator) + path)
		// This can't fail, as (by definition) all paths are relative to root.
		path, _ = filepath.Rel(string(os.PathSeparator), path)
	}

	// Clean the path again for good measure.
	return filepath.Clean(path)
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
			if v, ok := client.GameState.Players.Get(info.Plid); ok {
				data, err := json.Marshal(playerState{Plid: info.Plid, State: v})
				if err != nil {
					continue
				}

				s.SendMessage(
					"/api/live",
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
			"/api/live",
			sse.NewMessage(
				"", // id
				string(data),
				"player-left", // event
			),
		)
	})

	c.On(func(client *session.InsimSession, mso *protocol.Mso) {
		if player, ok := c.GameState.Players.Get(mso.Plid); ok {
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
				"/api/live",
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
			fmt.Printf("%s conns=%d\n", info.HName, info.NumConns)
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
			"/api/live",
			sse.NewMessage(
				"",              // id
				string(chatmsg), // event
				"chat",
			),
		)

		gamedata, err := client.GameState.MarshalJSON()
		if err != nil {
			return
		}

		s.SendMessage(
			"/api/live",
			sse.NewMessage(
				"",               // id
				string(gamedata), // event
				"state",
			),
		)

	})

	c.On(func(client *session.InsimSession, res *protocol.Res) {
		if player, ok := c.GameState.Players.Get(res.Plid); ok {
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
				"/api/live",
				sse.NewMessage(
					"", // id
					string(data),
					"chat", // event
				),
			)
		}
	})

	c.On(func(client *session.InsimSession, con *protocol.Con) {
		a, aok := c.GameState.Players.Get(con.A.Plid)
		b, bok := c.GameState.Players.Get(con.B.Plid)

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
				"/api/live",
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

	r.Mount("/api/live", s)

	r.Get("/api/connections", func(w http.ResponseWriter, r *http.Request) {
		render.JSON(w, r, &c.GameState.Connections)
	})

	r.Get("/api/players", func(w http.ResponseWriter, r *http.Request) {
		render.JSON(w, r, &c.GameState.Players)
	})

	r.Get("/api/state", func(w http.ResponseWriter, r *http.Request) {
		render.JSON(w, r, &c.GameState)
	})

	r.Get("/track/{code}", func(w http.ResponseWriter, r *http.Request) {
		track := CleanPath(chi.URLParam(r, "code"))

		pth := files.Pth{}
		pth.Read(
			filepath.Join(
				"..",
				"pth",
				"data",
				fmt.Sprintf("%s.pth", track),
			),
		)

		fit := pth.FitTo(1024, 1024, 2)

		var b bytes.Buffer
		buf := io.Writer(&b)

		// trackcolour string, limitcolour string, linecolour string, startfinishcolour string

		const trackcolour = "#1F2937"
		const limitcolour = "#059669"
		const linecolour = "#F9FAFB"
		const startfinishcolour = "#ffffff"

		s := svg.New(buf)
		s.Start(1024, 1024, "style=\"border: 1px solid red\"")
		s.Polygon(
			fit.OuterX, fit.OuterY,
			fmt.Sprintf("stroke: %s; stroke-width:2px; fill: %s; fill-rule: evenodd", limitcolour, limitcolour),
		)
		s.Polygon(
			fit.RoadX, fit.RoadY,
			fmt.Sprintf("stroke: %s; stroke-width:2px; fill: %s; fill-rule: evenodd", trackcolour, trackcolour),
		)
		s.Line(
			fit.FinishX[0], fit.FinishY[0], fit.FinishX[1], fit.FinishY[1],
			fmt.Sprintf("stroke: %s; stroke-width: 2px;", startfinishcolour),
		)

		s.Polyline(
			fit.RoadCX, fit.RoadCY,
			fmt.Sprintf("stroke: %s; stroke-width: 2px; fill: none", linecolour),
		)
		s.End()

		w.Header().Set("Content-Type", "image/svg+xml")
		w.Write(
			[]byte(b.String()),
		)
	})

	r.Get("/api/track/{code}", func(w http.ResponseWriter, r *http.Request) {
		track := CleanPath(chi.URLParam(r, "code"))

		pth := files.Pth{}
		pth.Read(
			filepath.Join(
				"..",
				"pth",
				"data",
				fmt.Sprintf("%s.pth", track),
			),
		)

		fit := pth.FitTo(1024, 1024, 2)

		var payload struct {
			Fit   files.PthFit
			Image string
		}

		payload.Fit = fit
		payload.Image = fmt.Sprintf("/track/%s", track)

		render.JSON(w, r, payload)
	})

	http.ListenAndServe(":4000", r)
}
