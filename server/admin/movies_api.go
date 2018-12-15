package admin

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/cavaliercoder/grab"
	"github.com/h2non/filetype"
	"github.com/labstack/echo"
	"github.com/mragiadakos/borinema/server/conf"
	"github.com/mragiadakos/borinema/server/db"
	uuid "github.com/satori/go.uuid"
)

func (aa *adminApi) startDownloadAndUpdateDB(folder, url, uuid string, wsSend func(*db.DbMovie)) {
	dbmovie, err := db.GetMovieByUuid(aa.db, uuid)
	if err != nil {
		log.Println("Error:", err)
		dbmovie.State = db.MOVIE_STATE_ERROR
		dbmovie.Error = err.Error()
		dbmovie.Update(aa.db)
		return
	}

	client := grab.NewClient()
	req, err := grab.NewRequest(folder+"/"+uuid, url)
	if err != nil {
		log.Println("Error:", err)
		dbmovie.State = db.MOVIE_STATE_ERROR
		dbmovie.Error = err.Error()
		dbmovie.Update(aa.db)
		return
	}
	log.Printf("Info: Downloading %v...\n", req.URL())
	resp := client.Do(req)
	if resp.Err() != nil {
		dbmovie.State = db.MOVIE_STATE_ERROR
		dbmovie.Error = resp.Err().Error()
		log.Println("Error:", dbmovie.Error)
		dbmovie.Update(aa.db)
		return
	}
	log.Println(resp.HTTPResponse)
	log.Printf("Info: http status  %v\n", resp.HTTPResponse.Status)

	if resp.HTTPResponse.StatusCode >= http.StatusBadRequest {
		dbmovie.State = db.MOVIE_STATE_ERROR
		dbmovie.Error = "The link failed with status " + fmt.Sprint(resp.HTTPResponse.StatusCode)
		log.Println("Error:", dbmovie.Error)
		dbmovie.Update(aa.db)
		return
	}
	t := *time.NewTicker(500 * time.Millisecond)
	log.Println("starting the loop")
Loop:
	for {
		select {
		case <-t.C:
			dbmovie.Progress = 100 * resp.Progress()
			dbmovie.Update(aa.db)
			wsSend(dbmovie)
			log.Println("tick")

		case <-resp.Done:
			dbmovie.Progress = 100
			dbmovie.Update(aa.db)
			wsSend(dbmovie)
			break Loop
		}
	}

	t.Stop()
	buf, err := ioutil.ReadFile(folder + "/" + uuid)
	if err != nil {
		dbmovie.State = db.MOVIE_STATE_ERROR
		dbmovie.Error = "Can not read the file " + err.Error()
		dbmovie.Filetype = db.FILE_TYPE_OTHER
		log.Println("Error:", dbmovie.Error)
		dbmovie.Update(aa.db)
		wsSend(dbmovie)
		return
	}

	kind, err := filetype.Match(buf)
	if err != nil {
		dbmovie.State = db.MOVIE_STATE_ERROR
		dbmovie.Error = "The file is type of unknown " + err.Error()
		dbmovie.Filetype = db.FILE_TYPE_OTHER
		log.Println("Error:", dbmovie.Error)
		dbmovie.Update(aa.db)
		wsSend(dbmovie)
		return
	}

	log.Printf("Info: File type: %s. MIME: %s\n", kind.Extension, kind.MIME.Value)
	if kind.Extension != "mp4" && kind.Extension != "webm" {
		dbmovie.State = db.MOVIE_STATE_ERROR
		dbmovie.Error = "The file is type is different than mp4 and webm "
		dbmovie.Filetype = db.FILE_TYPE_OTHER
		log.Println("Error:", dbmovie.Error)
		dbmovie.Update(aa.db)
		wsSend(dbmovie)
		return
	}

	if kind.Extension == "mp4" {
		dbmovie.Filetype = db.FILE_TYPE_MP4
	}

	if kind.Extension == "webm" {
		dbmovie.Filetype = db.FILE_TYPE_WEBM
	}

	dbmovie.Progress = 100
	if err := resp.Err(); err != nil {
		dbmovie.State = db.MOVIE_STATE_ERROR
		dbmovie.Error = "The link failed with error " + err.Error()
		log.Println("Error:", dbmovie.Error)
		dbmovie.Update(aa.db)
		wsSend(dbmovie)
		return
	}
	dbmovie.State = db.MOVIE_STATE_FINISHED

	dbmovie.Update(aa.db)
	wsSend(dbmovie)

	log.Printf("Info: Download saved to %v \n", resp.Filename)
}

func (aa *adminApi) DownloadMovieLink(config conf.Configuration, wsSend func(*db.DbMovie)) func(c echo.Context) error {
	return func(c echo.Context) error {
		input := MovieFromLinkInput{}
		c.Bind(&input)
		al := AdminLogic{}

		createDbEntry := func(url, name string) (string, error) {
			id := uuid.NewV4()
			movie := &db.DbMovie{}
			movie.ID = id.String()
			movie.Name = name
			movie.Link = url
			movie.State = db.MOVIE_STATE_DOWNLOADING
			err := movie.Create(aa.db)
			return movie.ID, err
		}
		startDownload := func(url, id string) {
			go aa.startDownloadAndUpdateDB(config.DownloadFolder, url, id, wsSend)
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
	gmo.ID = dm.ID
	gmo.Name = dm.Name
	gmo.CreatedAt = dm.CreatedAt.UnixNano()
	gmo.Progress = dm.Progress
	gmo.State = string(dm.State)
	gmo.Filetype = string(dm.Filetype)
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
		limitStr := c.QueryParam("limit")
		lastSeenStr := c.QueryParam("last_seen_date")

		pagination := Pagination{}
		pagination.Limit, _ = strconv.Atoi(limitStr)
		if pagination.Limit == 0 {
			pagination.Limit = -1
		}
		lastSeen, _ := strconv.Atoi(lastSeenStr)
		if lastSeen > 0 {
			t := time.Unix(0, int64(lastSeen))
			pagination.LastSeenDate = &t
		}
		log.Println(pagination)
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
		al := AdminLogic{}
		errMsg := al.DeleteMovie(uuid, aa.movieExists, aa.deleteMovie(config))
		if errMsg != nil {
			return c.JSON(errMsg.Status, errMsg.Json())
		}
		return c.JSON(http.StatusNoContent, "")
	}
}

func (aa *adminApi) UpdateMovie(config conf.Configuration) func(c echo.Context) error {
	return func(c echo.Context) error {
		uuid := c.Param("id")
		input := UpdateMovieInput{}
		c.Bind(&input)
		al := AdminLogic{}
		errMsg := al.UpdateMovie(uuid, input, aa.movieExists, aa.updateMovie)
		if errMsg != nil {
			return c.JSON(errMsg.Status, errMsg.Json())
		}
		return c.JSON(http.StatusNoContent, "")
	}
}

func (aa *adminApi) SelectMovie(config conf.Configuration) func(c echo.Context) error {
	return func(c echo.Context) error {
		uuid := c.Param("id")
		al := AdminLogic{}

		errMsg := al.SelectMovie(uuid, aa.movieExists, aa.selectMovie)
		if errMsg != nil {
			return c.JSON(errMsg.Status, errMsg.Json())
		}
		return c.JSON(http.StatusNoContent, "")
	}
}

func (aa *adminApi) SelectedMovie(config conf.Configuration) func(c echo.Context) error {
	return func(c echo.Context) error {
		getSelectedMovieDb := func() (*MovieOutput, error) {
			dm, err := db.GetMovieBySelected(aa.db)
			if err != nil {
				return nil, err
			}
			m := aa.serializeMovie(*dm)
			return &m, nil
		}
		al := AdminLogic{}
		gmo, errMsg := al.SelectedMovie(getSelectedMovieDb)
		if errMsg != nil {
			return c.JSON(errMsg.Status, errMsg.Json())
		}
		return c.JSON(http.StatusOK, gmo)
	}
}

func (aa *adminApi) RemoveAnySelectedMovie(config conf.Configuration) func(c echo.Context) error {
	return func(c echo.Context) error {
		removeSelectedMovieDb := func() error {
			dm, err := db.GetMovieBySelected(aa.db)
			if err != nil {
				return err
			}
			dm.Selected = false
			return dm.Update(aa.db)
		}
		al := AdminLogic{}
		errMsg := al.RemoveAnySelectedMovie(removeSelectedMovieDb)
		if errMsg != nil {
			return c.JSON(errMsg.Status, errMsg.Json())
		}
		return c.JSON(http.StatusNoContent, "")
	}
}

func (aa *adminApi) movieExists(id string) bool {
	_, err := db.GetMovieByUuid(aa.db, id)
	return err == nil
}

func (aa *adminApi) updateMovie(id, name string) error {
	dbm, err := db.GetMovieByUuid(aa.db, id)
	if err != nil {
		return err
	}
	dbm.Name = name

	return dbm.Update(aa.db)
}

func (aa *adminApi) selectMovie(id string) error {
	dbms, _ := db.GetMoviesByPage(aa.db, -1, nil)
	for _, v := range dbms {
		v.Selected = false
		v.Update(aa.db)
	}
	dbm, err := db.GetMovieByUuid(aa.db, id)
	if err != nil {
		return err
	}
	dbm.Selected = true
	return dbm.Update(aa.db)
}

func (aa *adminApi) deleteMovie(config conf.Configuration) func(id string) error {
	return func(id string) error {
		dbm, err := db.GetMovieByUuid(aa.db, id)
		if err != nil {
			return err
		}
		os.Remove(config.DownloadFolder + "/" + id)
		return dbm.Delete(aa.db)
	}
}
