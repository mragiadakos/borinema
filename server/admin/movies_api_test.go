package admin

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/labstack/echo"
	"github.com/mragiadakos/borinema/server/conf"
	"github.com/stretchr/testify/assert"
)

func TestDownloadMovieLinkApiSuccess(t *testing.T) {
	config := conf.Configuration{}
	config.Folder = "/tmp/downloads"
	os.Remove("/tmp/test.db")
	db, err := NewDB("/tmp/test.db")
	assert.Nil(t, err)
	os.MkdirAll(config.Folder, 0777)

	e := echo.New()
	input := MovieFromLinkInput{
		Url: "https://sample-videos.com/video123/mp4/720/big_buck_bunny_720p_1mb.mp4",
	}
	b, _ := json.Marshal(input)
	req := httptest.NewRequest(http.MethodPost, "/api/admin/movies/link", strings.NewReader(string(b)))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	aa := NewAdminApi(db)
	aa.DownloadMovieLink(config)(c)
	assert.Equal(t, http.StatusOK, rec.Code)

	// check if the movie exists in the DB
	movie := DbMovie{}
	db.Model(&DbMovie{}).Last(&movie)
	assert.Equal(t, input.Url, movie.Link)
	assert.Equal(t, MovieDownloading, movie.State)

	time.Sleep(time.Second * 20)
	db.Model(&DbMovie{}).Last(&movie)
	assert.Equal(t, MovieFinished, movie.State)
	assert.Equal(t, FiletypeMp4, movie.Filetype)

	// check if the movie exists in the folder
	_, err = os.Stat(config.Folder + "/" + movie.Id)
	assert.Nil(t, err)
}

func TestGetMovieApiSuccess(t *testing.T) {
	config := conf.Configuration{}
	config.Folder = "/tmp/downloads"
	os.Remove("/tmp/test.db")
	db, err := NewDB("/tmp/test.db")
	assert.Nil(t, err)
	os.MkdirAll(config.Folder, 0777)

	// Download Movie
	e := echo.New()
	input := MovieFromLinkInput{
		Url: "https://sample-videos.com/video123/mp4/720/big_buck_bunny_720p_1mb.mp4",
	}
	b, _ := json.Marshal(input)
	req := httptest.NewRequest(http.MethodPost, "/api/admin/movies/link", strings.NewReader(string(b)))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	aa := NewAdminApi(db)
	aa.DownloadMovieLink(config)(c)
	output := MovieFromLinkOutput{}
	json.Unmarshal(rec.Body.Bytes(), &output)

	// Get movie in the state of downloading
	req = httptest.NewRequest(http.MethodGet, "/api/admin/movies/"+output.ID, nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec = httptest.NewRecorder()
	c = e.NewContext(req, rec)
	c.SetPath("/api/admin/movies/:id")
	c.SetParamNames("id")
	c.SetParamValues(output.ID)
	aa.GeMovie(config)(c)
	mout := &GetMovieOutput{}
	json.Unmarshal(rec.Body.Bytes(), &mout)
	assert.Equal(t, output.ID, mout.ID)
	assert.Equal(t, string(MovieDownloading), mout.State)

	time.Sleep(20 * time.Second)

	// Get movie in the state of downloading
	req = httptest.NewRequest(http.MethodGet, "/api/admin/movies/"+output.ID, nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec = httptest.NewRecorder()
	c = e.NewContext(req, rec)
	c.SetPath("/api/admin/movies/:id")
	c.SetParamNames("id")
	c.SetParamValues(output.ID)
	aa.GeMovie(config)(c)
	mout = &GetMovieOutput{}
	json.Unmarshal(rec.Body.Bytes(), &mout)
	assert.Equal(t, output.ID, mout.ID)
	assert.Equal(t, string(MovieFinished), mout.State)
	assert.Equal(t, string(FiletypeMp4), mout.Filetype)
}
