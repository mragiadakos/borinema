package utils

type WsTheme string

const (
	WS_THEME_DOWNLOAD_PROGRESS_MOVIE = WsTheme("download_progress_movie")
	WS_THEME_PLAYER_ACTION           = WsTheme("player_action")
)

type WsData struct {
	Theme WsTheme     `json:"theme"`
	Data  interface{} `json:"data"`
}

type WsProgressMovieJson struct {
	ID       string  `json:"id"`
	State    string  `json:"state"`
	Progress float64 `json:"progress"`
	Filetype string  `json:"file_type"`
}
