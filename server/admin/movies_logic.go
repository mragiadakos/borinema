package admin

import (
	"errors"
	"net/http"
	"regexp"
	"time"

	"github.com/mragiadakos/borinema/server/utils"
)

type MovieFromLinkInput struct {
	Url  string `json:"url"`
	Name string `json:"name"`
}

type MovieFromLinkOutput struct {
	ID string `json:"id"`
}

func (al AdminLogic) DownloadMovieFromLink(
	input MovieFromLinkInput,
	createEntryOnDB func(url, name string) (id string, err error),
	startDownload func(url string, id string),
) (*MovieFromLinkOutput, *utils.ErrorMsg) {
	validUrl := regexp.MustCompile(`^(?:http(s)?:\/\/)?[\w.-]+(?:\.[\w\.-]+)+[\w\-\._~:/?#[\]@!\$&'\(\)\*\+,;=.]+$`)
	errMsg := utils.NewErrorMsg()
	if !validUrl.MatchString(input.Url) {
		errMsg.VariableErrors["url"] = ERR_URL_NOT_CORRECT
		errMsg.Status = http.StatusUnprocessableEntity
	}
	if len(input.Name) == 0 {
		errMsg.VariableErrors["name"] = ERR_NAME_IS_EMPTY
		errMsg.Status = http.StatusUnprocessableEntity
	}
	if errMsg.HasErrors() {
		errStr := ""
		for _, v := range errMsg.VariableErrors {
			errStr += v.Error() + "\n"
		}
		errMsg.Error = errors.New(errStr)
		return nil, errMsg
	}

	output := &MovieFromLinkOutput{}
	var err error
	output.ID, err = createEntryOnDB(input.Url, input.Name)
	if err != nil {
		errMsg.Error = ERR_DB_FAILED(err)
		errMsg.Status = http.StatusInternalServerError
		return nil, errMsg
	}
	startDownload(input.Url, output.ID)
	return output, nil
}

type MovieOutput struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Progress  float64   `json:"progress"`
	State     string    `json:"state"`
	Filetype  string    `json:"filetype"`
	CreatedAt time.Time `json:"created_at"`
	Error     string    `json:"error"`
}

func (al AdminLogic) GetMovie(id string,
	getDbMovie func(id string) (*MovieOutput, error)) (*MovieOutput, *utils.ErrorMsg) {
	errMsg := &utils.ErrorMsg{}
	gmo, err := getDbMovie(id)
	if err != nil {
		errMsg.Status = http.StatusNotFound
		errMsg.Error = ERR_MOVIE_NOT_FOUND
		return nil, errMsg
	}
	return gmo, nil
}

type Pagination struct {
	LastSeenDate *time.Time `json:"last_seen_date"`
	Limit        int        `json:"limit"` // -1 means all
}

func (al AdminLogic) GetMovies(pag Pagination, getDbMovies func(Pagination) []MovieOutput) ([]MovieOutput, *utils.ErrorMsg) {
	errMsg := utils.NewErrorMsg()
	if pag.Limit == 0 {
		errMsg.VariableErrors["limit"] = ERR_ITEMS_NOT_ZERO
		errMsg.Status = http.StatusUnprocessableEntity
	} else if pag.Limit < -1 {
		errMsg.VariableErrors["limit"] = ERR_ITEMS_NOT_LESS_MINUS_ONE
		errMsg.Status = http.StatusUnprocessableEntity
	}
	if errMsg.HasErrors() {
		return nil, errMsg
	}
	output := getDbMovies(pag)
	return output, nil
}

func (al AdminLogic) DeleteMovie(id string,
	movieExist func(id string) bool,
	deleteDb func(id string) error) *utils.ErrorMsg {
	errMsg := utils.NewErrorMsg()
	if !movieExist(id) {
		errMsg.Error = ERR_MOVIE_NOT_FOUND
		errMsg.Status = http.StatusNotFound
		return errMsg
	}
	err := deleteDb(id)
	if err != nil {
		errMsg.Error = err
		errMsg.Status = http.StatusInternalServerError
		return errMsg
	}
	return nil
}

type UpdateMovieInput struct {
	Name string `json:"name"`
}

func (al AdminLogic) UpdateMovie(id string,
	input UpdateMovieInput,
	movieExist func(id string) bool,
	updateMovie func(id, name string) error) *utils.ErrorMsg {
	errMsg := utils.NewErrorMsg()
	if !movieExist(id) {
		errMsg.Error = ERR_MOVIE_NOT_FOUND
		errMsg.Status = http.StatusNotFound
		return errMsg
	}
	if len(input.Name) == 0 {
		errMsg.VariableErrors["name"] = ERR_NAME_IS_EMPTY
		errMsg.Status = http.StatusUnprocessableEntity
		return errMsg
	}

	err := updateMovie(id, input.Name)
	if err != nil {
		errMsg.Error = err
		errMsg.Status = http.StatusInternalServerError
		return errMsg
	}
	return nil
}

func (al AdminLogic) SelectMovie(
	id string,
	movieExist func(id string) bool,
	selectMovie func(id string) error) *utils.ErrorMsg {
	errMsg := utils.NewErrorMsg()
	if !movieExist(id) {
		errMsg.Error = ERR_MOVIE_NOT_FOUND
		errMsg.Status = http.StatusNotFound
		return errMsg
	}

	err := selectMovie(id)
	if err != nil {
		errMsg.Error = err
		errMsg.Status = http.StatusInternalServerError
		return errMsg
	}
	return nil
}

func (al AdminLogic) SelectedMovie(
	selectedMovie func() (*MovieOutput, error)) (*MovieOutput, *utils.ErrorMsg) {
	errMsg := &utils.ErrorMsg{}
	gmo, err := selectedMovie()
	if err != nil {
		errMsg.Status = http.StatusNotFound
		errMsg.Error = ERR_MOVIE_NOT_FOUND
		return nil, errMsg
	}
	return gmo, nil
}

func (al AdminLogic) RemoveAnySelectedMovie(
	removeSelection func() error) *utils.ErrorMsg {
	errMsg := &utils.ErrorMsg{}
	err := removeSelection()
	if err != nil {
		errMsg.Status = http.StatusNotFound
		errMsg.Error = ERR_MOVIE_NOT_FOUND
		return errMsg
	}
	return nil
}
