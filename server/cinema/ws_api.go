package cinema

import (
	"encoding/json"
	"log"

	"github.com/labstack/echo"
	"github.com/mragiadakos/borinema/server/utils"
	"github.com/olahol/melody"
)

type cinemaWsApi struct {
	m     *melody.Melody
	users []*melody.Session
}

func NewCinemaWsApi(cp *utils.CommonMoviePlayer) *cinemaWsApi {
	a := &cinemaWsApi{}
	a.m = melody.New()
	a.users = []*melody.Session{}
	go cp.Receiver(a.onPlayerAction)
	return a
}

func (wa *cinemaWsApi) HttpFunc() func(echo.Context) error {
	return func(c echo.Context) error {
		wa.m.HandleConnect(wa.onConnection)
		wa.m.HandleRequest(c.Response(), c.Request())
		return nil
	}
}

func (wa *cinemaWsApi) onConnection(s *melody.Session) {
	log.Println("added new websocket session to users")
	wa.users = append(wa.users, s)
}

func (wa *cinemaWsApi) onPlayerAction(a utils.MoviePlayerAction) {
	for _, v := range wa.users {
		wd := utils.WsData{
			Theme: utils.WS_THEME_PLAYER_ACTION,
			Data:  a,
		}
		b, _ := json.Marshal(wd)
		v.Write(b)
	}
}
