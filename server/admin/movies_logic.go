package admin

import (
	"net/http"
	"regexp"

	"github.com/mragiadakos/borinema/server/utils"
)

type MovieFromLinkInput struct {
	Url string `json:"url"`
}

type MovieFromLinkOutput struct {
	ID string `json:"id"`
}

func (al AdminLogic) DownloadMovieFromLink(
	input MovieFromLinkInput,
	createEntryOnDB func(url string) (id string, err error),
	startDownload func(url string, id string),
) (*MovieFromLinkOutput, *utils.ErrorMsg) {
	validUrl := regexp.MustCompile(`^(?:http(s)?:\/\/)?[\w.-]+(?:\.[\w\.-]+)+[\w\-\._~:/?#[\]@!\$&'\(\)\*\+,;=.]+$`)
	errMsg := utils.NewErrorMsg()
	if !validUrl.MatchString(input.Url) {
		errMsg.VariableErrors["url"] = ERR_URL_NOT_CORRECT
		errMsg.Status = http.StatusUnprocessableEntity
	}
	if errMsg.HasErrors() {
		return nil, errMsg
	}

	output := &MovieFromLinkOutput{}
	var err error
	output.ID, err = createEntryOnDB(input.Url)
	if err != nil {
		errMsg.Error = ERR_DB_FAILED(err)
		errMsg.Status = http.StatusInternalServerError
		return nil, errMsg
	}
	startDownload(input.Url, output.ID)
	return output, nil
}

type GetMovieOutput struct {
	ID       string  `json:"id"`
	Progress float64 `json:"progress"`
	State    string  `json:"state"`
	Filetype string  `json:"filetype"`
	Error    string  `json:"error"`
}

func (al AdminLogic) GetMovie(id string,
	getDbMovie func(id string) (*GetMovieOutput, error)) (*GetMovieOutput, *utils.ErrorMsg) {
	errMsg := &utils.ErrorMsg{}
	gmo, err := getDbMovie(id)
	if err != nil {
		errMsg.Status = http.StatusNotFound
		errMsg.Error = ERR_MOVIE_NOT_FOUND
		return nil, errMsg
	}
	return gmo, nil
}
