package admin

import (
	"errors"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDownloadMovieLinkFail(t *testing.T) {
	aa := AdminLogic{}
	input := MovieFromLinkInput{Url: "hththththtlas://lallaslsad.com"}
	createEntry := func(url string) (string, error) {
		return "none", nil
	}
	downloader := func(s, id string) {}
	_, errMsg := aa.DownloadMovieFromLink(input, createEntry, downloader)
	assert.NotNil(t, errMsg)
	assert.Equal(t, errMsg.VariableErrors["url"], ERR_URL_NOT_CORRECT)
	assert.Equal(t, errMsg.Status, http.StatusUnprocessableEntity)
}

func TestDownloadMovieLinkSuccess(t *testing.T) {
	al := AdminLogic{}
	input := MovieFromLinkInput{Url: "http://lallaslsad.com/movie"}
	createEntry := func(url string) (string, error) {
		return "1234", nil
	}
	downloader := func(s, id string) {}
	out, errMsg := al.DownloadMovieFromLink(input, createEntry, downloader)
	assert.Nil(t, errMsg)
	assert.Equal(t, out.ID, "1234")
}

func TestGetMovieFailure(t *testing.T) {
	al := AdminLogic{}
	getMovieDb := func(id string) (*GetMovieOutput, error) {
		return nil, errors.New("failed")
	}
	_, errMsg := al.GetMovie("aaaa", getMovieDb)
	assert.Equal(t, http.StatusNotFound, errMsg.Status)
}
