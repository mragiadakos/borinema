package cinema

import (
	"net/http"

	"github.com/jinzhu/gorm"
	"github.com/labstack/echo"
	"github.com/mragiadakos/borinema/server/conf"
	"github.com/mragiadakos/borinema/server/db"
	"github.com/mragiadakos/borinema/server/utils"
)

type movieApi struct {
	db                 *gorm.DB
	config             conf.Configuration
	requestCurrentTime func()
}

func NewCinemaMovieApi(db *gorm.DB, config conf.Configuration, requestCurrentTimeWs func()) *movieApi {
	ma := &movieApi{}
	ma.db = db
	ma.config = config
	ma.requestCurrentTime = requestCurrentTimeWs
	return ma
}

func (ma *movieApi) GetMovieVideo() func(c echo.Context) error {
	return func(c echo.Context) error {
		uuid := c.Param("id")
		dbm, err := db.GetMovieByUuid(ma.db, uuid)
		if err != nil {
			errMsg := utils.NewErrorMsg()
			errMsg.Status = http.StatusNotFound
			errMsg.Error = ERR_MOVIE_NOT_FOUND
			return c.JSON(errMsg.Status, errMsg.Json())
		}
		return c.File(ma.config.DownloadFolder + "/" + dbm.ID)
	}
}

type MovieInfoOutput struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Filetype string `json:"filetype"`
}

func (ma *movieApi) serializeMovieInfo(dm db.DbMovie) MovieInfoOutput {
	gmo := MovieInfoOutput{}
	gmo.Name = dm.Name
	gmo.ID = dm.ID
	gmo.Filetype = string(dm.Filetype)
	return gmo
}

func (ma *movieApi) GetMovieInfo() func(c echo.Context) error {
	return func(c echo.Context) error {
		dbm, err := db.GetMovieBySelected(ma.db)
		if err != nil {
			errMsg := utils.NewErrorMsg()
			errMsg.Status = http.StatusNotFound
			errMsg.Error = ERR_MOVIE_NOT_FOUND
			return c.JSON(errMsg.Status, errMsg.Json())
		}

		mio := ma.serializeMovieInfo(*dbm)
		return c.JSON(http.StatusOK, mio)
	}
}

func (ma *movieApi) RequestCurrentTime() func(c echo.Context) error {
	return func(c echo.Context) error {
		_, err := db.GetMovieBySelected(ma.db)
		if err != nil {
			return c.JSON(http.StatusUnprocessableEntity, nil)
		}
		ma.requestCurrentTime()
		return c.JSON(http.StatusOK, nil)
	}
}

func (ma *movieApi) CinemaPage() func(c echo.Context) error {
	return func(c echo.Context) error {
		return c.File("cinema_panel/index.html")
	}
}
