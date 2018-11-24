package main

import (
	"strings"

	"github.com/cathalgarvey/fmtless/encoding/json"
	"github.com/oskca/gopherjs-localStorage"
	"honnef.co/go/js/xhr"
)

type ErrorMsg struct {
	VariableErrors map[string]string `json:"variable_errors"`
	Error          string            `json:"error"`
}

type AuthorizationJson struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
type TokenJson struct {
	Token string `json:"token"`
}

type AuthService struct{}

func (as AuthService) postLogin(auth AuthorizationJson) (string, *ErrorMsg) {
	b, err := json.Marshal(auth)
	if err != nil {
		return "", &ErrorMsg{Error: err.Error()}
	}
	req := xhr.NewRequest("POST", "/api/login")
	req.SetRequestHeader("Content-Type", "application/json")
	err = req.Send(b)
	if err != nil {
		return "", &ErrorMsg{Error: err.Error()}
	}
	text := req.ResponseText
	data := []byte(text)
	if strings.Contains(text, "error") {
		errMsg := ErrorMsg{}
		err = json.Unmarshal(data, &errMsg)
		if err != nil {
			return "", &ErrorMsg{Error: err.Error()}
		}
		return "", &errMsg
	} else {
		token := TokenJson{}
		err = json.Unmarshal(data, &token)
		if err != nil {
			return "", &ErrorMsg{Error: err.Error()}
		}
		return token.Token, nil
	}
}

func (as AuthService) saveToken(token string) {
	localStorage.SetItem("token", token)
}

func (as AuthService) getToken() string {
	return localStorage.GetItem("token")
}
