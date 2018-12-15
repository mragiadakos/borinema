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
	config.DownloadFolder = "/tmp/downloads"
	os.Remove("/tmp/test.db")
	dbtx, err := db.NewDB("/tmp/test.db")
	assert.Nil(t, err)
	os.MkdirAll(config.DownloadFolder, 0777)

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
	fakeWs := func(a string, c float64) {}
	aa.DownloadMovieLink(config, fakeWs)(c)
	assert.Equal(t, http.StatusOK, rec.Code)

	// check if the movie exists in the DB
	movie := db.DbMovie{}
	dbtx.Model(&db.DbMovie{}).Last(&movie)
	assert.Equal(t, input.Url, movie.Link)
	assert.Equal(t, db.MOVIE_STATE_DOWNLOADING, movie.State)
	assert.Equal(t, input.Name, movie.Name)

	time.Sleep(time.Second * 20)
	dbtx.Model(&db.DbMovie{}).Last(&movie)
	assert.Equal(t, db.MOVIE_STATE_FINISHED, movie.State)
	assert.Equal(t, db.FILE_TYPE_MP4, movie.Filetype)

	// check if the movie exists in the folder
	_, err = os.Stat(config.DownloadFolder + "/" + movie.ID)
	assert.Nil(t, err)
}

func TestGetMovieApiSuccess(t *testing.T) {
	config := conf.Configuration{}
	config.DownloadFolder = "/tmp/downloads"
	os.Remove("/tmp/test.db")
	dbtx, err := db.NewDB("/tmp/test.db")
	assert.Nil(t, err)
	os.MkdirAll(config.DownloadFolder, 0777)

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
	fakeWs := func(a string, c float64) {}
	aa.DownloadMovieLink(config, fakeWs)(c)
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
	assert.Equal(t, string(db.MOVIE_STATE_DOWNLOADING), mout.State)

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
	assert.Equal(t, string(db.MOVIE_STATE_FINISHED), mout.State)
	assert.Equal(t, string(db.FILE_TYPE_MP4), mout.Filetype)
}

func TestGetMoviesSuccess(t *testing.T) {
	config := conf.Configuration{}
	config.DownloadFolder = "/tmp/downloads"
	os.Remove("/tmp/test.db")
	dbtx, err := db.NewDB("/tmp/test.db")
	assert.Nil(t, err)
	os.MkdirAll(config.DownloadFolder, 0777)

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
	assert.Equal(t, ms[9].ID, mouts[0].ID)
	assert.Equal(t, ms[8].ID, mouts[1].ID)

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
	assert.Equal(t, ms[9].ID, mouts[0].ID)
	assert.Equal(t, ms[8].ID, mouts[1].ID)

}

func TestDeleteMovieSuccess(t *testing.T) {
	config := conf.Configuration{}
	config.DownloadFolder = "/tmp/downloads"
	os.Remove("/tmp/test.db")
	dbtx, err := db.NewDB("/tmp/test.db")
	assert.Nil(t, err)
	os.MkdirAll(config.DownloadFolder, 0777)

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
		c.SetParamValues(v.ID)
		aa.DeleteMovie(config)(c)
		assert.Equal(t, http.StatusNoContent, rec.Code)
	}
	mouts, _ := db.GetMoviesByPage(dbtx, -1, nil)
	assert.Equal(t, 2, len(mouts))
}

func TestUpdateMovieSuccess(t *testing.T) {
	config := conf.Configuration{}
	config.DownloadFolder = "/tmp/downloads"
	os.Remove("/tmp/test.db")
	dbtx, err := db.NewDB("/tmp/test.db")
	assert.Nil(t, err)
	os.MkdirAll(config.DownloadFolder, 0777)
	oldName := "Buck bunny"
	newName := "Duck funny"

	// add new movie
	e := echo.New()
	input := MovieFromLinkInput{
		Url:  "https://ia800703.us.archive.org/34/items/1mbFile/1mb.mp4",
		Name: oldName,
	}
	b, _ := json.Marshal(input)
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(string(b)))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	fakeWs := func(a string, c float64) {}
	aa := NewAdminApi(dbtx)
	aa.DownloadMovieLink(config, fakeWs)(c)
	assert.Equal(t, http.StatusOK, rec.Code)
	output := MovieFromLinkOutput{}
	json.Unmarshal(rec.Body.Bytes(), &output)

	// update movie
	input2 := UpdateMovieInput{Name: newName}
	b, _ = json.Marshal(input2)
	req = httptest.NewRequest(http.MethodPut, "/", strings.NewReader(string(b)))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec = httptest.NewRecorder()
	c = e.NewContext(req, rec)
	c.SetPath("/api/admin/movies/:id")
	c.SetParamNames("id")
	c.SetParamValues(output.ID)
	aa.UpdateMovie(config)(c)
	assert.Equal(t, http.StatusNoContent, rec.Code)

	// get movie
	dbm, err := db.GetMovieByUuid(dbtx, output.ID)
	assert.NoError(t, err)
	assert.Equal(t, newName, dbm.Name)
}

func TestSelectMovieSuccess(t *testing.T) {
	config := conf.Configuration{}
	config.DownloadFolder = "/tmp/downloads"
	os.Remove("/tmp/test.db")
	dbtx, err := db.NewDB("/tmp/test.db")
	assert.Nil(t, err)
	os.MkdirAll(config.DownloadFolder, 0777)
	e := echo.New()
	aa := NewAdminApi(dbtx)

	uuids := []string{}
	// add two movies
	for i := 0; i < 2; i++ {
		input := MovieFromLinkInput{
			Url:  "https://ia800703.us.archive.org/34/items/1mbFile/1mb.mp4",
			Name: "oldName",
		}
		b, _ := json.Marshal(input)
		req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(string(b)))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		fakeWs := func(a string, c float64) {}
		aa.DownloadMovieLink(config, fakeWs)(c)
		assert.Equal(t, http.StatusOK, rec.Code)
		output := MovieFromLinkOutput{}
		json.Unmarshal(rec.Body.Bytes(), &output)
		uuids = append(uuids, output.ID)
	}

	// select first id
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/api/admin/movies/:id/select")
	c.SetParamNames("id")
	c.SetParamValues(uuids[0])
	aa.SelectMovie(config)(c)
	assert.Equal(t, http.StatusNoContent, rec.Code)

	dbm, err := db.GetMovieByUuid(dbtx, uuids[0])
	assert.NoError(t, err)
	assert.True(t, dbm.Selected)
	dbm, err = db.GetMovieByUuid(dbtx, uuids[1])
	assert.NoError(t, err)
	assert.False(t, dbm.Selected)

	//select second
	req = httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec = httptest.NewRecorder()
	c = e.NewContext(req, rec)
	c.SetPath("/api/admin/movies/:id/select")
	c.SetParamNames("id")
	c.SetParamValues(uuids[1])
	aa.SelectMovie(config)(c)
	assert.Equal(t, http.StatusNoContent, rec.Code)

	dbm, err = db.GetMovieByUuid(dbtx, uuids[1])
	assert.NoError(t, err)
	assert.True(t, dbm.Selected)
	dbm, err = db.GetMovieByUuid(dbtx, uuids[0])
	assert.NoError(t, err)
	assert.False(t, dbm.Selected)

	// get selected movie
	req = httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec = httptest.NewRecorder()
	c = e.NewContext(req, rec)
	c.SetPath("/api/admin/movies/selected")
	aa.SelectedMovie(config)(c)
	assert.Equal(t, http.StatusOK, rec.Code)

	mout := &MovieOutput{}
	json.Unmarshal(rec.Body.Bytes(), &mout)
	assert.Equal(t, uuids[1], mout.ID)

	// remove selections
	req = httptest.NewRequest(http.MethodDelete, "/", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec = httptest.NewRecorder()
	c = e.NewContext(req, rec)
	c.SetPath("/api/admin/movies/selected")
	aa.RemoveAnySelectedMovie(config)(c)
	assert.Equal(t, http.StatusNoContent, rec.Code)

	// get selected movie, but dont find any
	req = httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec = httptest.NewRecorder()
	c = e.NewContext(req, rec)
	c.SetPath("/api/admin/movies/selected")
	aa.SelectedMovie(config)(c)
	assert.Equal(t, http.StatusNotFound, rec.Code)

	// check that all are false
	dbms, err := db.GetMovies(dbtx)
	assert.NoError(t, err)
	assert.Equal(t, 2, len(dbms))
	for _, dbm := range dbms {
		assert.False(t, dbm.Selected)
	}

}
