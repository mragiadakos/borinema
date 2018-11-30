package services

import (
	"encoding/json"
	"time"

	"honnef.co/go/js/xhr"
)

type MovieService struct{}

type AddMovieJson struct {
	Url  string `json:"url"`
	Name string `json:"name"`
}

type MovieIdJson struct {
	ID string `json:"id"`
}

func (ms MovieService) AddMovie(amj AddMovieJson) (*MovieIdJson, *ErrorMsg) {
	as := AuthService{}
	token := as.GetToken()
	req := xhr.NewRequest("POST", "/api/admin/movies/link")
	req.SetRequestHeader("Content-Type", "application/json")
	req.SetRequestHeader("Authorization", "Bearer "+token)
	b, _ := json.Marshal(amj)
	err := req.Send(b)
	if err != nil {
		return nil, &ErrorMsg{Error: err.Error()}
	}
	data := []byte(req.ResponseText)
	if IsErrorMsg(data) {
		errMsg := ErrorMsg{}
		err = json.Unmarshal(data, &errMsg)
		if err != nil {
			return nil, &ErrorMsg{Error: err.Error()}
		}
		return nil, &errMsg
	}
	mij := MovieIdJson{}
	err = json.Unmarshal(data, &mij)
	if err != nil {
		return nil, &ErrorMsg{Error: err.Error()}
	}
	return &mij, nil
}

type PaginationJson struct {
	LastSeenDate *time.Time `json:"last_seen_date"`
	Limit        int        `json:"limit"` // -1 means all
}
type MovieJson struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Progress  float64   `json:"progress"`
	State     string    `json:"state"`
	Filetype  string    `json:"filetype"`
	CreatedAt time.Time `json:"created_at"`
	Error     string    `json:"error"`
}

func (ms MovieService) GetMovies(pag PaginationJson) ([]MovieJson, *ErrorMsg) {
	as := AuthService{}
	token := as.GetToken()
	req := xhr.NewRequest("POST", "/api/admin/get/movies")
	req.SetRequestHeader("Content-Type", "application/json")
	req.SetRequestHeader("Authorization", "Bearer "+token)
	b, _ := json.Marshal(pag)
	err := req.Send(b)
	if err != nil {
		return nil, &ErrorMsg{Error: err.Error()}
	}
	data := []byte(req.ResponseText)
	if IsErrorMsg(data) {
		errMsg := ErrorMsg{}
		err = json.Unmarshal(data, &errMsg)
		if err != nil {
			return nil, &ErrorMsg{Error: err.Error()}
		}
		return nil, &errMsg
	}
	mjs := []MovieJson{}
	err = json.Unmarshal(data, &mjs)
	if err != nil {
		return nil, &ErrorMsg{Error: err.Error()}
	}
	return mjs, nil
}

func (ms MovieService) GetMovie(id string) (*MovieJson, *ErrorMsg) {
	as := AuthService{}
	token := as.GetToken()
	req := xhr.NewRequest("GET", "/api/admin/movies/"+id)
	req.SetRequestHeader("Content-Type", "application/json")
	req.SetRequestHeader("Authorization", "Bearer "+token)
	err := req.Send(nil)
	if err != nil {
		return nil, &ErrorMsg{Error: err.Error()}
	}
	data := []byte(req.ResponseText)
	if IsErrorMsg(data) {
		errMsg := ErrorMsg{}
		err = json.Unmarshal(data, &errMsg)
		if err != nil {
			return nil, &ErrorMsg{Error: err.Error()}
		}
		return nil, &errMsg
	}
	mj := &MovieJson{}
	err = json.Unmarshal(data, &mj)
	if err != nil {
		return nil, &ErrorMsg{Error: err.Error()}
	}
	return mj, nil
}

func (ms MovieService) DeleteMovie(id string) *ErrorMsg {
	as := AuthService{}
	token := as.GetToken()
	req := xhr.NewRequest("DELETe", "/api/admin/movies/"+id)
	req.SetRequestHeader("Content-Type", "application/json")
	req.SetRequestHeader("Authorization", "Bearer "+token)
	err := req.Send(nil)
	if err != nil {
		return &ErrorMsg{Error: err.Error()}
	}
	data := []byte(req.ResponseText)
	if IsErrorMsg(data) {
		errMsg := ErrorMsg{}
		err = json.Unmarshal(data, &errMsg)
		if err != nil {
			return &ErrorMsg{Error: err.Error()}
		}
		return &errMsg
	}
	return nil
}

func (ms MovieService) RemoveMovieSelection() *ErrorMsg {
	as := AuthService{}
	token := as.GetToken()
	req := xhr.NewRequest("DELETe", "/api/admin/movies/selected")
	req.SetRequestHeader("Content-Type", "application/json")
	req.SetRequestHeader("Authorization", "Bearer "+token)
	err := req.Send(nil)
	if err != nil {
		return &ErrorMsg{Error: err.Error()}
	}
	data := []byte(req.ResponseText)
	if IsErrorMsg(data) {
		errMsg := ErrorMsg{}
		err = json.Unmarshal(data, &errMsg)
		if err != nil {
			return &ErrorMsg{Error: err.Error()}
		}
		return &errMsg
	}
	return nil
}

func (ms MovieService) SelectMovie(id string) *ErrorMsg {
	as := AuthService{}
	token := as.GetToken()
	req := xhr.NewRequest("DELETe", "/api/admin/movies/"+id+"/select")
	req.SetRequestHeader("Content-Type", "application/json")
	req.SetRequestHeader("Authorization", "Bearer "+token)
	err := req.Send(nil)
	if err != nil {
		return &ErrorMsg{Error: err.Error()}
	}
	data := []byte(req.ResponseText)
	if IsErrorMsg(data) {
		errMsg := ErrorMsg{}
		err = json.Unmarshal(data, &errMsg)
		if err != nil {
			return &ErrorMsg{Error: err.Error()}
		}
		return &errMsg
	}
	return nil
}
