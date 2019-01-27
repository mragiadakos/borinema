package services

import (
	"errors"

	"github.com/gopherjs/gopherjs/js"

	"github.com/gopherjs/websocket/websocketjs"
	"honnef.co/go/js/dom"
)

type MoviePlayerActionType string

type MoviePlayerAction struct {
	Action    MoviePlayerActionType `json:"action"`
	Time      float64               `json:"time"`
	IsPlaying bool                  `json:"is_playing"`
}

const (
	PLAY         = MoviePlayerActionType("play")
	STOP         = MoviePlayerActionType("stop")
	PAUSE        = MoviePlayerActionType("pause")
	CURRENT_TIME = MoviePlayerActionType("current_time")
)

var PlayerWs *websocketjs.WebSocket

type WsService struct{}

func init() {
	host := dom.GetWindow().Location().Host
	playerWs, err := websocketjs.New("ws://" + host + "/api/player/ws")
	if err != nil {
		println("Error for player WS :" + err.Error())
		return
	}
	PlayerWs = playerWs

}

const (
	WS_THEME_PLAYER_ACTION = "player_action"
)

type WsData struct {
	Theme string      `json:"theme"`
	Data  interface{} `json:"data"`
}

func (wss *WsService) SerializeMovieAction(data WsData) (*MoviePlayerAction, error) {
	jm, ok := data.Data.(map[string]interface{})
	if !ok {
		return nil, errors.New("The 'data' is not type of map[string]interface ")
	}

	mpa := MoviePlayerAction{}
	action, ok := jm["action"].(string)
	if !ok {
		return nil, errors.New("The 'data' is missing the key 'action' ")
	}
	t, ok := jm["time"].(float64)
	if !ok {
		return nil, errors.New("The 'data' is missing the key 'time' ")
	}
	isPlaying, ok := jm["is_playing"].(bool)
	if !ok {
		return nil, errors.New("The 'data' is missing the key 'is_playing' ")
	}
	mpa.Action = MoviePlayerActionType(action)
	mpa.Time = t
	mpa.IsPlaying = isPlaying
	return &mpa, nil
}

func (wss *WsService) SetPlayerWsOnMessage(f func(ev *js.Object)) {
	PlayerWs.AddEventListener("message", false, f)
}

func (wss *WsService) SetPlayerWsOnOpen(f func(ev *js.Object)) {
	PlayerWs.AddEventListener("open", false, f)
}

func (wss *WsService) SetPlayerWsOnClose(f func(ev *js.Object)) {
	PlayerWs.AddEventListener("close", false, f)
}
