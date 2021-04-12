// +build ignore

package main

import (
	"flag"
	"fmt"
	"time"

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
		for {
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
				fmt.Println(err)
			}

			time.Sleep(3 * time.Second)
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

	e.GET("/", func(ec echo.Context) error {
		return ec.HTML(
			http.StatusOK,
			`<!DOCTYPE html>
<html lang="en">
<head>
<title>Chat Example</title>
<script type="text/javascript">
window.onload = function () {
    var conn;
    var msg = document.getElementById("msg");
    var log = document.getElementById("log");

    function appendLog(item) {
        var doScroll = log.scrollTop > log.scrollHeight - log.clientHeight - 1;
        log.appendChild(item);
        if (doScroll) {
            log.scrollTop = log.scrollHeight - log.clientHeight;
        }
    }

    document.getElementById("form").onsubmit = function () {
        if (!conn) {
            return false;
        }
        if (!msg.value) {
            return false;
        }
        conn.send(msg.value);
        msg.value = "";
        return false;
    };

    if (window["WebSocket"]) {
        conn = new WebSocket("ws://" + document.location.host + "/ws");
        conn.onclose = function (evt) {
            var item = document.createElement("div");
            item.innerHTML = "<b>Connection closed.</b>";
            appendLog(item);
        };
        conn.onmessage = function (evt) {
					console.log(evt)
        };
    } else {
        var item = document.createElement("div");
        item.innerHTML = "<b>Your browser does not support WebSockets.</b>";
        appendLog(item);
    }
};
</script>
<style type="text/css">
html {
    overflow: hidden;
}

body {
    overflow: hidden;
    padding: 0;
    margin: 0;
    width: 100%;
    height: 100%;
    background: gray;
}

#log {
    background: white;
    margin: 0;
    padding: 0.5em 0.5em 0.5em 0.5em;
    position: absolute;
    top: 0.5em;
    left: 0.5em;
    right: 0.5em;
    bottom: 3em;
    overflow: auto;
}

#form {
    padding: 0 0.5em 0 0.5em;
    margin: 0;
    position: absolute;
    bottom: 1em;
    left: 0px;
    width: 100%;
    overflow: hidden;
}

</style>
</head>
<body>
<div id="log"></div>
<form id="form">
    <input type="submit" value="Send" />
    <input type="text" id="msg" size="64" autofocus />
</form>
</body>
</html>`,
		)
	})

  e.Logger.Fatal(e.Start("0.0.0.0:4000"))
}
