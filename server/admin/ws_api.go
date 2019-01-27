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
	m            *melody.Melody
	admins       []*melody.Session
	commonPlayer *utils.CommonMoviePlayer
	config       conf.Configuration
}

func NewAdminWsApi(config conf.Configuration, commonPlayer *utils.CommonMoviePlayer) *adminWsApi {
	a := &adminWsApi{}
	a.m = melody.New()
	a.admins = []*melody.Session{}
	a.commonPlayer = commonPlayer
	a.config = config
	return a
}

func (wa *adminWsApi) HttpFunc() func(echo.Context) error {
	return func(c echo.Context) error {
		tokenString := c.QueryParam("token")
		isAdmin, errMsg := utils.ReadAdminTokenString(wa.config, tokenString)
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
		wa.m.HandleMessage(wa.onMessage)
		wa.m.HandleRequest(c.Response(), c.Request())
		return nil
	}
}

func (wa *adminWsApi) onConnection(s *melody.Session) {
	log.Println("added new websocket session to admins")
	wa.admins = append(wa.admins, s)
}

func (wa *adminWsApi) onMessage(s *melody.Session, msg []byte) {
	wd := utils.WsData{}
	json.Unmarshal(msg, &wd)
	log.Println("Received from admin's WS :" + string(msg))
	switch wd.Theme {
	case utils.WS_THEME_PLAYER_ACTION:
		b, _ := json.Marshal(wd.Data)
		a := utils.MoviePlayerAction{}
		err := json.Unmarshal(b, &a)
		if err != nil {
			log.Println("Error: the message is not action.")
			return
		}
		wa.commonPlayer.Sender(a)
	}
}

func (wa *adminWsApi) SendProgressOfMovie(dbm *db.DbMovie) {
	log.Println("send to ", wa.admins)
	wp := utils.WsProgressMovieJson{}
	wp.ID = dbm.ID
	wp.Progress = dbm.Progress
	wp.State = dbm.State.String()
	wp.Filetype = dbm.Filetype.String()
	wd := utils.WsData{
		Theme: utils.WS_THEME_DOWNLOAD_PROGRESS_MOVIE,
		Data:  wp,
	}
	b, _ := json.Marshal(wd)
	for _, v := range wa.admins {

		v.Write(b)
	}
}

func (wa *adminWsApi) RequestCurrentTime() {
	a := utils.MoviePlayerAction{}
	a.Action = utils.REQUEST_CURRENT_TIME
	wd := utils.WsData{
		Theme: utils.WS_THEME_PLAYER_ACTION,
		Data:  a,
	}
	b, _ := json.Marshal(wd)
	for _, v := range wa.admins {
		v.Write(b)
	}
}
