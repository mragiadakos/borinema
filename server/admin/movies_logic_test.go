package admin

import (
	"errors"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDownloadMovieLinkFailOnUrl(t *testing.T) {
	aa := AdminLogic{}
	input := MovieFromLinkInput{Url: "hththththtlas://lallaslsad.com", Name: "lala"}
	createEntry := func(url, name string) (string, error) {
		return "none", nil
	}
	downloader := func(s, id string) {}
	_, errMsg := aa.DownloadMovieFromLink(input, createEntry, downloader)
	assert.NotNil(t, errMsg)
	assert.Equal(t, errMsg.VariableErrors["url"], ERR_URL_NOT_CORRECT)
	assert.Equal(t, errMsg.Status, http.StatusUnprocessableEntity)
}
func TestDownloadMovieLinkFailOnEmptyName(t *testing.T) {
	aa := AdminLogic{}
	input := MovieFromLinkInput{Url: "http://lallaslsad.com", Name: ""}
	createEntry := func(url, name string) (string, error) {
		return "none", nil
	}
	downloader := func(s, id string) {}
	_, errMsg := aa.DownloadMovieFromLink(input, createEntry, downloader)
	assert.NotNil(t, errMsg)
	assert.Equal(t, errMsg.VariableErrors["name"], ERR_NAME_IS_EMPTY)
	assert.Equal(t, errMsg.Status, http.StatusUnprocessableEntity)
}

func TestDownloadMovieLinkSuccess(t *testing.T) {
	al := AdminLogic{}
	input := MovieFromLinkInput{Url: "http://lallaslsad.com/movie", Name: "name"}
	createEntry := func(url, name string) (string, error) {
		return "1234", nil
	}
	downloader := func(s, id string) {}
	out, errMsg := al.DownloadMovieFromLink(input, createEntry, downloader)
	assert.Nil(t, errMsg)
	assert.Equal(t, out.ID, "1234")
}

func TestGetMovieFailure(t *testing.T) {
	al := AdminLogic{}
	getMovieDb := func(id string) (*MovieOutput, error) {
		return nil, errors.New("failed")
	}
	_, errMsg := al.GetMovie("aaaa", getMovieDb)
	assert.Equal(t, http.StatusNotFound, errMsg.Status)
}

func TestGetMoviesFailureOnPagination(t *testing.T) {
	al := AdminLogic{}
	pag := Pagination{
		Limit: 0,
	}
	_, errMsg := al.GetMovies(pag, func(Pagination) []MovieOutput { return nil })
	assert.Equal(t, http.StatusUnprocessableEntity, errMsg.Status)
	assert.Equal(t, errMsg.VariableErrors["limit"], ERR_ITEMS_NOT_ZERO)

	pag = Pagination{
		Limit: -2,
	}
	_, errMsg = al.GetMovies(pag, func(Pagination) []MovieOutput { return nil })
	assert.Equal(t, http.StatusUnprocessableEntity, errMsg.Status)
	assert.Equal(t, errMsg.VariableErrors["limit"], ERR_ITEMS_NOT_LESS_MINUS_ONE)
}

func TestDeleteMovieFailureNotFound(t *testing.T) {
	al := AdminLogic{}
	movieExists := func(id string) bool {
		return false
	}
	deleteMovie := func(id string) error {
		return nil
	}
	errMsg := al.DeleteMovie("lalalla", movieExists, deleteMovie)
	assert.Equal(t, errMsg.Error, ERR_MOVIE_NOT_FOUND)
}

func TestUpdateMovieFailureNotFound(t *testing.T) {
	al := AdminLogic{}
	movieExists := func(id string) bool {
		return false
	}
	updateMovie := func(id, name string) error {
		return nil
	}
	input := UpdateMovieInput{
		Name: "lalal",
	}
	errMsg := al.UpdateMovie("lalalla", input, movieExists, updateMovie)
	assert.Equal(t, errMsg.Error, ERR_MOVIE_NOT_FOUND)
}

func TestUpdateMovieFailureEmptyName(t *testing.T) {
	al := AdminLogic{}
	movieExists := func(id string) bool {
		return true
	}
	updateMovie := func(id, name string) error {
		return nil
	}
	input := UpdateMovieInput{
		Name: "",
	}
	errMsg := al.UpdateMovie("lalalla", input, movieExists, updateMovie)
	assert.Equal(t, errMsg.VariableErrors["name"], ERR_NAME_IS_EMPTY)
}

func TestSelectMovieFailureNotFound(t *testing.T) {
	al := AdminLogic{}
	movieExists := func(id string) bool {
		return false
	}
	selectMovie := func(id string) error {
		return nil
	}
	errMsg := al.SelectMovie("lalalla", movieExists, selectMovie)
	assert.Equal(t, errMsg.Error, ERR_MOVIE_NOT_FOUND)
}

func TestSelectedMovieFailure(t *testing.T) {
	al := AdminLogic{}
	getMovieDb := func() (*MovieOutput, error) {
		return nil, errors.New("failed")
	}
	_, errMsg := al.SelectedMovie(getMovieDb)
	assert.Equal(t, http.StatusNotFound, errMsg.Status)
}
