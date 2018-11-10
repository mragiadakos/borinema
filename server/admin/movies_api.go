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
	uuid "github.com/satori/go.uuid"
)

func (aa *adminApi) startDownloadAndUpdateDB(folder, url, uuid string) {
	movie := DbMovie{}
	err := aa.db.Model(&DbMovie{}).Last(&movie).Error
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
	log.Printf("Info: http status  %v\n", resp.HTTPResponse.Status)
	if resp.HTTPResponse.StatusCode >= http.StatusBadRequest {
		movie.State = MovieError
		movie.Error = "The link failed with status " + fmt.Sprint(resp.HTTPResponse.StatusCode)
		movie.Update(aa.db)
		return
	}
	t := time.NewTicker(500 * time.Millisecond)
	defer t.Stop()

Loop:
	for {
		select {
		case <-t.C:
			movie.Progress = 100 * resp.Progress()
			movie.Update(aa.db)

		case <-resp.Done:
			break Loop
		}
	}

	buf, err := ioutil.ReadFile(folder + "/" + uuid)
	if err != nil {
		movie.State = MovieError
		movie.Error = "Can not read the file " + err.Error()
		movie.Filetype = FiletypeOther
		movie.Update(aa.db)
	}

	kind, err := filetype.Match(buf)
	if err != nil {
		movie.State = MovieError
		movie.Error = "The file is type of unknown " + err.Error()
		movie.Filetype = FiletypeOther
		movie.Update(aa.db)
		return
	}

	log.Printf("Info: File type: %s. MIME: %s\n", kind.Extension, kind.MIME.Value)
	if kind.Extension != "mp4" && kind.Extension != "webm" {
		movie.State = MovieError
		movie.Error = "The file is type is different than mp4 and webm "
		movie.Filetype = FiletypeOther
		movie.Update(aa.db)
		return
	}

	if kind.Extension == "mp4" {
		movie.Filetype = FiletypeMp4
	}

	if kind.Extension == "webm" {
		movie.Filetype = FiletypeWebm
	}

	movie.Progress = 100
	if err := resp.Err(); err != nil {
		movie.State = MovieError
		movie.Error = "The link failed with error " + err.Error()
		movie.Update(aa.db)
		return
	}
	movie.State = MovieFinished

	movie.Update(aa.db)

	log.Printf("Info: Download saved to %v \n", resp.Filename)
}

func (aa *adminApi) DownloadMovieLink(config conf.Configuration) func(c echo.Context) error {
	return func(c echo.Context) error {
		input := MovieFromLinkInput{}
		c.Bind(&input)
		al := AdminLogic{}

		createDbEntry := func(url string) (string, error) {
			id := uuid.NewV4()
			movie := &DbMovie{}
			movie.Id = id.String()
			movie.Link = url
			movie.State = MovieDownloading
			err := movie.Create(aa.db)
			return movie.Id, err
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

func (aa *adminApi) GeMovie(config conf.Configuration) func(c echo.Context) error {
	return func(c echo.Context) error {
		id := c.Param("id")
		log.Println("ID get movie " + id)
		getMovieDb := func(id string) (*GetMovieOutput, error) {
			dm := &DbMovie{}
			err := aa.db.Model(&DbMovie{}).Where("id = ?", id).Find(&dm).Error
			if err != nil {
				return nil, err
			}
			gmo := &GetMovieOutput{}
			gmo.ID = dm.Id
			gmo.Progress = dm.Progress
			gmo.State = string(dm.State)
			gmo.Filetype = string(dm.Filetype)
			gmo.Error = dm.Error
			return gmo, nil
		}
		al := AdminLogic{}
		gmo, errMsg := al.GetMovie(id, getMovieDb)
		if errMsg != nil {
			return c.JSON(errMsg.Status, errMsg.Json())
		}
		return c.JSON(http.StatusOK, gmo)
	}
}
