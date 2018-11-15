package admin

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/labstack/echo"
	"github.com/mragiadakos/borinema/server/conf"
	"github.com/mragiadakos/borinema/server/db"
	"github.com/stretchr/testify/assert"
)

func TestDownloadMovieLinkApiSuccess(t *testing.T) {
	config := conf.Configuration{}
	config.Folder = "/tmp/downloads"
	os.Remove("/tmp/test.db")
	dbtx, err := db.NewDB("/tmp/test.db")
	assert.Nil(t, err)
	os.MkdirAll(config.Folder, 0777)

	e := echo.New()
	input := MovieFromLinkInput{
		Url:  "https://ia800703.us.archive.org/34/items/1mbFile/1mb.mp4",
		Name: "Buck bunny",
	}
	b, _ := json.Marshal(input)
	req := httptest.NewRequest(http.MethodPost, "/api/admin/movies/link", strings.NewReader(string(b)))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	aa := NewAdminApi(dbtx)
	aa.DownloadMovieLink(config)(c)
	assert.Equal(t, http.StatusOK, rec.Code)

	// check if the movie exists in the DB
	movie := db.DbMovie{}
	dbtx.Model(&db.DbMovie{}).Last(&movie)
	assert.Equal(t, input.Url, movie.Link)
	assert.Equal(t, db.MovieDownloading, movie.State)
	assert.Equal(t, input.Name, movie.Name)

	time.Sleep(time.Second * 20)
	dbtx.Model(&db.DbMovie{}).Last(&movie)
	assert.Equal(t, db.MovieFinished, movie.State)
	assert.Equal(t, db.FiletypeMp4, movie.Filetype)

	// check if the movie exists in the folder
	_, err = os.Stat(config.Folder + "/" + movie.Uuid)
	assert.Nil(t, err)
}

func TestGetMovieApiSuccess(t *testing.T) {
	config := conf.Configuration{}
	config.Folder = "/tmp/downloads"
	os.Remove("/tmp/test.db")
	dbtx, err := db.NewDB("/tmp/test.db")
	assert.Nil(t, err)
	os.MkdirAll(config.Folder, 0777)

	// Download Movie
	e := echo.New()
	input := MovieFromLinkInput{
		Url:  "https://ia800703.us.archive.org/34/items/1mbFile/1mb.mp4",
		Name: "Buck bunny",
	}
	b, _ := json.Marshal(input)
	req := httptest.NewRequest(http.MethodPost, "/api/admin/movies/link", strings.NewReader(string(b)))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	aa := NewAdminApi(dbtx)
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
	aa.GetMovie(config)(c)
	mout := &MovieOutput{}
	json.Unmarshal(rec.Body.Bytes(), &mout)
	assert.Equal(t, output.ID, mout.ID)
	assert.Equal(t, string(db.MovieDownloading), mout.State)

	time.Sleep(20 * time.Second)

	// Get movie in the state of downloading
	req = httptest.NewRequest(http.MethodGet, "/api/admin/movies/"+output.ID, nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec = httptest.NewRecorder()
	c = e.NewContext(req, rec)
	c.SetPath("/api/admin/movies/:id")
	c.SetParamNames("id")
	c.SetParamValues(output.ID)
	aa.GetMovie(config)(c)
	mout = &MovieOutput{}
	json.Unmarshal(rec.Body.Bytes(), &mout)
	assert.Equal(t, output.ID, mout.ID)
	assert.Equal(t, input.Name, mout.Name)
	assert.Equal(t, string(db.MovieFinished), mout.State)
	assert.Equal(t, string(db.FiletypeMp4), mout.Filetype)
}

func TestGetMoviesSuccess(t *testing.T) {
	config := conf.Configuration{}
	config.Folder = "/tmp/downloads"
	os.Remove("/tmp/test.db")
	dbtx, err := db.NewDB("/tmp/test.db")
	assert.Nil(t, err)
	os.MkdirAll(config.Folder, 0777)

	ms := []db.DbMovie{}
	for i := 0; i < 10; i++ {
		m := db.DbMovie{}
		m.Name = "Movie " + fmt.Sprint(i)
		err := m.Create(dbtx)
		assert.Nil(t, err)
		time.Sleep(1 * time.Second)
		ms = append(ms, m)
	}

	aa := NewAdminApi(dbtx)
	input := Pagination{
		Limit: 2,
	}
	b, _ := json.Marshal(input)
	req := httptest.NewRequest(http.MethodGet, "/api/admin/movies", strings.NewReader(string(b)))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	e := echo.New()
	c := e.NewContext(req, rec)
	aa.GetMovies(config)(c)
	mouts := []MovieOutput{}
	json.Unmarshal(rec.Body.Bytes(), &mouts)
	assert.Equal(t, 2, len(mouts))
	assert.Equal(t, ms[9].Uuid, mouts[0].ID)
	assert.Equal(t, ms[8].Uuid, mouts[1].ID)

	input = Pagination{
		LastSeenDate: &ms[7].CreatedAt,
		Limit:        2,
	}
	b, _ = json.Marshal(input)
	req = httptest.NewRequest(http.MethodGet, "/api/admin/movies", strings.NewReader(string(b)))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec = httptest.NewRecorder()
	c = e.NewContext(req, rec)
	aa.GetMovies(config)(c)
	mouts = []MovieOutput{}
	json.Unmarshal(rec.Body.Bytes(), &mouts)
	assert.Equal(t, 2, len(mouts))
	assert.Equal(t, ms[9].Uuid, mouts[0].ID)
	assert.Equal(t, ms[8].Uuid, mouts[1].ID)

}

func TestDeleteMovieSuccess(t *testing.T) {
	config := conf.Configuration{}
	config.Folder = "/tmp/downloads"
	os.Remove("/tmp/test.db")
	dbtx, err := db.NewDB("/tmp/test.db")
	assert.Nil(t, err)
	os.MkdirAll(config.Folder, 0777)

	ms := []db.DbMovie{}
	for i := 0; i < 10; i++ {
		m := db.DbMovie{}
		m.Name = "Movie " + fmt.Sprint(i)
		err := m.Create(dbtx)
		assert.Nil(t, err)
		time.Sleep(1 * time.Second)
		ms = append(ms, m)
	}
	aa := NewAdminApi(dbtx)
	for _, v := range ms[2:] {
		req := httptest.NewRequest(http.MethodDelete, "/", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		e := echo.New()
		c := e.NewContext(req, rec)
		c.SetPath("/api/admin/movies/:id")
		c.SetParamNames("id")
		c.SetParamValues(v.Uuid)
		aa.DeleteMovie(config)(c)
		assert.Equal(t, http.StatusNoContent, rec.Code)
	}
	mouts, _ := db.GetMoviesByPage(dbtx, -1, nil)
	assert.Equal(t, 2, len(mouts))
}
