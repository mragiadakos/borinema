package main

import (
	"encoding/json"

	"github.com/gopherjs/gopherjs/js"
	"github.com/mragiadakos/borinema/admin_panel/actions"
	"github.com/mragiadakos/borinema/admin_panel/services"
	"github.com/mragiadakos/borinema/admin_panel/store"
)

func SetMovies() {
	go func() {
		ms := services.MovieService{}
		pag := services.PaginationJson{}
		pag.Limit = 10
		mvs, _ := ms.GetMovies(pag)
		store.Dispatch(&actions.SetMovies{
			Movies: mvs,
		})
	}()
}

func UpdateMovieProgress(wsm services.WsProgressMovieJson) {
	go func() {
		store.Dispatch(&actions.SetMovieProgress{
			ID:       wsm.Id,
			Progress: wsm.Progress,
			State:    wsm.State,
			Filetype: wsm.Filetype,
		})
	}()
}

func EnableWebsocket() {
	go func() {
		movieIsPlaying := false
		wss := &services.WsService{}
		wss.SetAdminWsOnMessage(func(ev *js.Object) {
			println("admin ws message: " + ev.Get("data").String())
			b := []byte(ev.Get("data").String())
			wsd := services.WsData{}
			json.Unmarshal(b, &wsd)
			switch wsd.Theme {
			case services.WS_THEME_DOWNLOAD_PROGRESS_MOVIE:
				wsm, err := wss.SerializeProgressMovie(wsd)
				if err != nil {
					println("admin ws error: " + err.Error())
					return
				}
				UpdateMovieProgress(*wsm)
			case services.WS_THEME_PLAYER_ACTION:
				mpa, err := wss.SerializeMovieAction(wsd)
				if err != nil {
					println("player ws error: " + err.Error())
					return
				}
				if mpa.Action == services.REQUEST_CURRENT_TIME {
					videoPlayer := js.Global.Get("document").Call("getElementById", "media-video")
					currentTime := videoPlayer.Get("currentTime").Float()
					wss.SendCurrentTimeToAdmin(services.CURRENT_TIME, currentTime, movieIsPlaying)
				}
			default:
				println("admin ws error: could not find the 'theme' of the ws message.")
			}
		})
		wss.SetAdminWsOnClose(func(ev *js.Object) {
			println("admin ws close: " + ev.Get("data").String())
		})
		wss.SetAdminWsOnOpen(func(ev *js.Object) {
			println("admin ws open: " + ev.Get("data").String())
		})
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
					movieIsPlaying = true
				case services.PAUSE:
					videoPlayer.Call("pause")
					movieIsPlaying = false
				case services.STOP:
					videoPlayer.Call("pause")
					videoPlayer.Set("currentTime", 0)
					movieIsPlaying = false
				case services.CURRENT_TIME:
					videoPlayer.Set("currentTime", mpa.Time)
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
