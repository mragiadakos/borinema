package admin

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/cavaliercoder/grab"
	"github.com/h2non/filetype"
	"github.com/labstack/echo"
	"github.com/mragiadakos/borinema/server/conf"
	"github.com/mragiadakos/borinema/server/db"
	uuid "github.com/satori/go.uuid"
)

func (aa *adminApi) startDownloadAndUpdateDB(folder, url, uuid string) {
	dbmovie, err := db.GetMovieByUuid(aa.db, uuid)
	if err != nil {
		log.Println("Error:", err)
		return
	}

	client := grab.NewClient()
	req, err := grab.NewRequest(folder+"/"+uuid, url)
	if err != nil {
		log.Println("Error:", err)
		return
	}
	log.Printf("Info: Downloading %v...\n", req.URL())
	resp := client.Do(req)
	log.Println(resp.HTTPResponse)
	log.Printf("Info: http status  %v\n", resp.HTTPResponse.Status)

	if resp.HTTPResponse.StatusCode >= http.StatusBadRequest {
		dbmovie.State = db.MovieError
		dbmovie.Error = "The link failed with status " + fmt.Sprint(resp.HTTPResponse.StatusCode)
		dbmovie.Update(aa.db)
		return
	}
	t := time.NewTicker(500 * time.Millisecond)
	defer t.Stop()

Loop:
	for {
		select {
		case <-t.C:
			dbmovie.Progress = 100 * resp.Progress()
			dbmovie.Update(aa.db)

		case <-resp.Done:
			break Loop
		}
	}

	buf, err := ioutil.ReadFile(folder + "/" + uuid)
	if err != nil {
		dbmovie.State = db.MovieError
		dbmovie.Error = "Can not read the file " + err.Error()
		dbmovie.Filetype = db.FiletypeOther
		dbmovie.Update(aa.db)
	}

	kind, err := filetype.Match(buf)
	if err != nil {
		dbmovie.State = db.MovieError
		dbmovie.Error = "The file is type of unknown " + err.Error()
		dbmovie.Filetype = db.FiletypeOther
		dbmovie.Update(aa.db)
		return
	}

	log.Printf("Info: File type: %s. MIME: %s\n", kind.Extension, kind.MIME.Value)
	if kind.Extension != "mp4" && kind.Extension != "webm" {
		dbmovie.State = db.MovieError
		dbmovie.Error = "The file is type is different than mp4 and webm "
		dbmovie.Filetype = db.FiletypeOther
		dbmovie.Update(aa.db)
		return
	}

	if kind.Extension == "mp4" {
		dbmovie.Filetype = db.FiletypeMp4
	}

	if kind.Extension == "webm" {
		dbmovie.Filetype = db.FiletypeWebm
	}

	dbmovie.Progress = 100
	if err := resp.Err(); err != nil {
		dbmovie.State = db.MovieError
		dbmovie.Error = "The link failed with error " + err.Error()
		dbmovie.Update(aa.db)
		return
	}
	dbmovie.State = db.MovieFinished

	dbmovie.Update(aa.db)

	log.Printf("Info: Download saved to %v \n", resp.Filename)
}

func (aa *adminApi) DownloadMovieLink(config conf.Configuration) func(c echo.Context) error {
	return func(c echo.Context) error {
		input := MovieFromLinkInput{}
		c.Bind(&input)
		al := AdminLogic{}

		createDbEntry := func(url, name string) (string, error) {
			id := uuid.NewV4()
			movie := &db.DbMovie{}
			movie.Uuid = id.String()
			movie.Name = name
			movie.Link = url
			movie.State = db.MovieDownloading
			err := movie.Create(aa.db)
			return movie.Uuid, err
		}
		startDownload := func(url, id string) {
			go aa.startDownloadAndUpdateDB(config.Folder, url, id)
		}
		output, errMsg := al.DownloadMovieFromLink(input, createDbEntry, startDownload)
		if errMsg != nil {
			return c.JSON(errMsg.Status, errMsg.Json())
		}
		return c.JSON(http.StatusOK, output)
	}
}
func (aa *adminApi) serializeMovie(dm db.DbMovie) MovieOutput {
	gmo := MovieOutput{}
	gmo.ID = dm.Uuid
	gmo.Name = dm.Name
	gmo.CreatedAt = dm.CreatedAt
	gmo.Progress = dm.Progress
	gmo.State = string(dm.State)
	gmo.Filetype = string(dm.Filetype)
	gmo.Error = dm.Error
	return gmo
}
func (aa *adminApi) GetMovie(config conf.Configuration) func(c echo.Context) error {
	return func(c echo.Context) error {
		uuid := c.Param("id")
		log.Println("ID get movie " + uuid)
		getMovieDb := func(uuid string) (*MovieOutput, error) {
			dm, err := db.GetMovieByUuid(aa.db, uuid)
			if err != nil {
				return nil, err
			}
			m := aa.serializeMovie(*dm)
			return &m, nil
		}
		al := AdminLogic{}
		gmo, errMsg := al.GetMovie(uuid, getMovieDb)
		if errMsg != nil {
			return c.JSON(errMsg.Status, errMsg.Json())
		}
		return c.JSON(http.StatusOK, gmo)
	}
}

func (aa *adminApi) GetMovies(config conf.Configuration) func(c echo.Context) error {
	return func(c echo.Context) error {
		al := AdminLogic{}
		pagination := Pagination{}
		c.Bind(&pagination)
		getMoviesDb := func(pagination Pagination) []MovieOutput {
			dbms, _ := db.GetMoviesByPage(aa.db, pagination.Limit, pagination.LastSeenDate)
			movies := []MovieOutput{}
			for _, v := range dbms {
				movies = append(movies, aa.serializeMovie(v))
			}
			return movies
		}
		movies, errMsg := al.GetMovies(pagination, getMoviesDb)
		if errMsg != nil {
			return c.JSON(errMsg.Status, errMsg.Json())
		}

		return c.JSON(http.StatusOK, movies)
	}
}

func (aa *adminApi) DeleteMovie(config conf.Configuration) func(c echo.Context) error {
	return func(c echo.Context) error {
		uuid := c.Param("id")
		movieExists := func(id string) bool {
			_, err := db.GetMovieByUuid(aa.db, id)
			return err == nil
		}
		deleteMovie := func(id string) error {
			dbm, err := db.GetMovieByUuid(aa.db, id)
			if err != nil {
				return err
			}
			return dbm.Delete(aa.db)
		}
		al := AdminLogic{}
		errMsg := al.DeleteMovie(uuid, movieExists, deleteMovie)
		if errMsg != nil {
			return c.JSON(errMsg.Status, errMsg.Json())
		}
		return c.JSON(http.StatusNoContent, "")
	}
}
