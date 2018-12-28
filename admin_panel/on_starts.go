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
		wss := &services.WsService{}
		as := services.AuthService{}
		wss, err := services.NewWsService(as.GetToken())
		if err != nil {
			println(err)
		}
		println(wss)
		wss.SetOnMessage(func(ev *js.Object) {
			println("ws message: " + ev.Get("data").String())
			b := []byte(ev.Get("data").String())
			wsd := services.WsData{}
			json.Unmarshal(b, &wsd)
			switch wsd.Theme {
			case services.WS_THEME_DOWNLOAD_PROGRESS_MOVIE:
				wsm, err := wss.SerializeProgressMovie(wsd)
				if err != nil {
					println("ws error: " + err.Error())
					return
				}
				UpdateMovieProgress(*wsm)
			default:
				println("ws error: could not find the 'theme' of the ws message.")
			}
		})
		wss.SetOnClose(func(ev *js.Object) {
			println("ws close: " + ev.Get("data").String())
		})
		wss.SetOnOpen(func(ev *js.Object) {
			println("ws open: " + ev.Get("data").String())
		})
	}()
}
