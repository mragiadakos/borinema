package services

import (
	"errors"

	"github.com/gopherjs/gopherjs/js"
	"github.com/gopherjs/websocket/websocketjs"
)

type WsService struct {
	ws *websocketjs.WebSocket
}

func NewWsService(token string) (*WsService, error) {
	wss := &WsService{}
	ws, err := websocketjs.New("ws://localhost:8080/api/admin/ws?token=" + token)
	if err != nil {
		return nil, err
	}
	wss.ws = ws

	return wss, nil
}

const (
	WS_THEME_PROGRESS_MOVIE = "progress_movie"
)

type WsData struct {
	Theme string      `json:"theme"`
	Data  interface{} `json:"data"`
}
type WsProgressMovieJson struct {
	Id       string  `json:"id"`
	State    string  `json:"state"`
	Progress float64 `json:"progress"`
	Filetype string  `json:"file_type"`
}

func (wss *WsService) SerializeProgressMovie(data WsData) (*WsProgressMovieJson, error) {
	jm, ok := data.Data.(map[string]interface{})
	if !ok {
		return nil, errors.New("The 'data' is not type of map[string]interface ")
	}

	wpm := WsProgressMovieJson{}
	wpm.Id, ok = jm["id"].(string)
	if !ok {
		return nil, errors.New("The 'data' is missing the key 'id' ")
	}
	wpm.Progress, ok = jm["progress"].(float64)
	if !ok {
		return nil, errors.New("The 'data' is missing the key 'progress' ")
	}
	wpm.State, ok = jm["state"].(string)
	if !ok {
		return nil, errors.New("The 'data' is missing the key 'state' ")
	}
	wpm.Filetype, ok = jm["file_type"].(string)
	if !ok {
		return nil, errors.New("The 'data' is missing the key 'file_type' ")
	}
	return &wpm, nil
}

func (wss *WsService) SetOnMessage(f func(ev *js.Object)) {
	wss.ws.AddEventListener("message", false, f)
}

func (wss *WsService) SetOnOpen(f func(ev *js.Object)) {
	wss.ws.AddEventListener("open", false, f)
}

func (wss *WsService) SetOnClose(f func(ev *js.Object)) {
	wss.ws.AddEventListener("close", false, f)
}
