package services

import (
	"encoding/json"
	"errors"

	"honnef.co/go/js/dom"

	"github.com/gopherjs/gopherjs/js"
	"github.com/gopherjs/websocket/websocketjs"
)

var AdminWs *websocketjs.WebSocket
var PlayerWs *websocketjs.WebSocket

type WsService struct{}

func init() {
	as := AuthService{}
	host := dom.GetWindow().Location().Host
	adminWs, err := websocketjs.New("ws://" + host + "/api/admin/ws?token=" + as.GetToken())
	if err != nil {
		println("Error for admin WS :" + err.Error())
		return
	}
	AdminWs = adminWs

	playerWs, err := websocketjs.New("ws://" + host + "/api/player/ws")
	if err != nil {
		println("Error for player WS :" + err.Error())
		return
	}
	PlayerWs = playerWs

}

const (
	WS_THEME_DOWNLOAD_PROGRESS_MOVIE = "download_progress_movie"
	WS_THEME_PLAYER_ACTION           = "player_action"
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

type MoviePlayerActionType string

type MoviePlayerAction struct {
	Action    MoviePlayerActionType `json:"action"`
	Time      float64               `json:"time"`
	IsPlaying bool                  `json:"is_playing"`
}

const (
	PLAY                 = MoviePlayerActionType("play")
	STOP                 = MoviePlayerActionType("stop")
	PAUSE                = MoviePlayerActionType("pause")
	REQUEST_CURRENT_TIME = MoviePlayerActionType("request_current_time")
	CURRENT_TIME         = MoviePlayerActionType("current_time")
)

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

func (wss *WsService) SetAdminWsOnMessage(f func(ev *js.Object)) {
	AdminWs.AddEventListener("message", false, f)
}

func (wss *WsService) SetPlayerWsOnMessage(f func(ev *js.Object)) {
	PlayerWs.AddEventListener("message", false, f)
}

func (wss *WsService) SetAdminWsOnOpen(f func(ev *js.Object)) {
	AdminWs.AddEventListener("open", false, f)
}

func (wss *WsService) SetPlayerWsOnOpen(f func(ev *js.Object)) {
	PlayerWs.AddEventListener("open", false, f)
}

func (wss *WsService) SetAdminWsOnClose(f func(ev *js.Object)) {
	AdminWs.AddEventListener("close", false, f)
}

func (wss *WsService) SetPlayerWsOnClose(f func(ev *js.Object)) {
	PlayerWs.AddEventListener("close", false, f)
}

func (wss *WsService) SendMoviePlayerActionToAdmin(action MoviePlayerActionType) {
	mpa := MoviePlayerAction{Action: action}
	wd := WsData{
		Theme: WS_THEME_PLAYER_ACTION,
		Data:  mpa,
	}
	b, _ := json.Marshal(wd)
	println("sended " + string(b))
	err := AdminWs.Send(string(b))
	if err != nil {
		println("Error: " + err.Error())
	}
}

func (wss *WsService) SendCurrentTimeToAdmin(action MoviePlayerActionType, t float64, isPlaying bool) {
	mpa := MoviePlayerAction{Action: action, Time: t, IsPlaying: isPlaying}
	wd := WsData{
		Theme: WS_THEME_PLAYER_ACTION,
		Data:  mpa,
	}
	b, _ := json.Marshal(wd)
	println("sended " + string(b))
	err := AdminWs.Send(string(b))
	if err != nil {
		println("Error: " + err.Error())
	}
}
