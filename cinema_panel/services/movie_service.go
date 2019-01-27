package services

import (
	"encoding/json"

	"honnef.co/go/js/xhr"
)

type MovieService struct{}

type ErrorMsg struct {
	VariableErrors map[string]string `json:"variable_errors"`
	Error          string            `json:"error"`
}

func IsErrorMsg(b []byte) bool {
	tmp := map[string]interface{}{}
	json.Unmarshal(b, &tmp)
	_, ok := tmp["error"]
	_, ok2 := tmp["variable_errors"]
	return ok || ok2
}

type MovieJson struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Filetype string `json:"filetype"`
}

func (ms MovieService) GetMovie() (*MovieJson, *ErrorMsg) {
	req := xhr.NewRequest("GET", "/api/cinema/movie")
	req.SetRequestHeader("Content-Type", "application/json")
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

func (ms MovieService) RequestCurrentTime() *ErrorMsg {
	req := xhr.NewRequest("POST", "/api/cinema/currentTime")
	req.SetRequestHeader("Content-Type", "application/json")
	err := req.Send(nil)
	if err != nil {
		return &ErrorMsg{Error: err.Error()}
	}
	if req.Status >= 400 {
		return &ErrorMsg{Error: "There is not a movie selected"}
	}
	return nil
}
