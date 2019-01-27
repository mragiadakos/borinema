package main

import (
	"encoding/json"

	"github.com/mragiadakos/borinema/cinema_panel/actions"
	"github.com/mragiadakos/borinema/cinema_panel/store"

	"github.com/gopherjs/gopherjs/js"
	"github.com/mragiadakos/borinema/cinema_panel/services"
)

func enableWebsocket() {
	go func() {
		wss := &services.WsService{}
		wss.SetPlayerWsOnMessage(func(ev *js.Object) {
			println("playerws message: " + ev.Get("data").String())
			b := []byte(ev.Get("data").String())
			wsd := services.WsData{}
			json.Unmarshal(b, &wsd)
			switch wsd.Theme {
			case services.WS_THEME_PLAYER_ACTION:
				mpa, err := wss.SerializeMovieAction(wsd)
				if err != nil {
					println("player ws error: " + err.Error())
					return
				}
				videoPlayer := js.Global.Get("document").Call("getElementById", "media-video")
				switch mpa.Action {
				case services.PLAY:
					videoPlayer.Call("play")
				case services.PAUSE:
					videoPlayer.Call("pause")
				case services.STOP:
					videoPlayer.Call("pause")
					videoPlayer.Set("currentTime", 0)
				case services.CURRENT_TIME:
					videoPlayer.Set("currentTime", mpa.Time)
					if mpa.IsPlaying {
						videoPlayer.Call("play")
					}
				default:
					println("player ws error: unknown player action")
				}
			default:
				println("player ws error: could not find the 'theme' of the ws message.")
			}

		})
		wss.SetPlayerWsOnClose(func(ev *js.Object) {
			println("player ws close: " + ev.Get("data").String())
		})
		wss.SetPlayerWsOnOpen(func(ev *js.Object) {
			println("player ws open: " + ev.Get("data").String())
		})

	}()
}

func getMovieOnStart() {
	go func() {
		ms := services.MovieService{}
		m, errj := ms.GetMovie()
		if errj != nil {
			store.Dispatch(&actions.SetMovie{})
			return
		}
		store.Dispatch(&actions.SetMovie{
			MovieID:  m.ID,
			Name:     m.Name,
			Filetype: m.Filetype,
		})
		go func() {
			ms.RequestCurrentTime()
		}()
	}()
}
