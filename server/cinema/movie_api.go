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
	db     *gorm.DB
	config conf.Configuration
}

func NewCinemaMovieApi(db *gorm.DB, config conf.Configuration) *movieApi {
	ma := &movieApi{}
	ma.db = db
	ma.config = config
	return ma
}

func (ma *movieApi) GetMovie() func(c echo.Context) error {
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
	Name string `json:"name"`
}

func (ma *movieApi) serializeMovieInfo(dm db.DbMovie) MovieInfoOutput {
	gmo := MovieInfoOutput{}
	gmo.Name = dm.Name
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
