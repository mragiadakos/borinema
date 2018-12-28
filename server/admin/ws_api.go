package admin

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/mragiadakos/borinema/server/db"

	"github.com/labstack/echo"
	"github.com/mragiadakos/borinema/server/conf"
	"github.com/mragiadakos/borinema/server/utils"
	"github.com/olahol/melody"
)

type adminWsApi struct {
	m      *melody.Melody
	admins []*melody.Session
}

func NewAdminWsApi() *adminWsApi {
	a := &adminWsApi{}
	a.m = melody.New()
	a.admins = []*melody.Session{}
	return a
}

func (wa *adminWsApi) HttpFunc(config conf.Configuration) func(echo.Context) error {
	return func(c echo.Context) error {
		tokenString := c.QueryParam("token")
		isAdmin, errMsg := utils.ReadAdminTokenString(config, tokenString)
		if errMsg != nil {
			return c.JSON(errMsg.Status, errMsg.Json())
		}
		if !isAdmin {
			errMsg := utils.NewErrorMsg()
			errMsg.Error = errors.New("You are not an admin.")
			errMsg.Status = http.StatusUnauthorized
			return c.JSON(errMsg.Status, errMsg.Json())
		}
		wa.m.HandleConnect(wa.onConnection)
		wa.m.HandleRequest(c.Response(), c.Request())
		return nil
	}
}

func (wa *adminWsApi) onConnection(s *melody.Session) {
	log.Println("added new websocket session to admins")
	wa.admins = append(wa.admins, s)
}

type WsTheme string

const (
	WS_THEME_DOWNLOAD_PROGRESS_MOVIE = WsTheme("download_progress_movie")
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

func (wa *adminWsApi) SendProgressOfMovie(dbm *db.DbMovie) {
	log.Println("send to ", wa.admins)
	for _, v := range wa.admins {
		wp := WsProgressMovieJson{}
		wp.ID = dbm.ID
		wp.Progress = dbm.Progress
		wp.State = dbm.State.String()
		wp.Filetype = dbm.Filetype.String()
		wd := WsData{
			Theme: WS_THEME_DOWNLOAD_PROGRESS_MOVIE,
			Data:  wp,
		}
		b, _ := json.Marshal(wd)
		v.Write(b)
	}
}
