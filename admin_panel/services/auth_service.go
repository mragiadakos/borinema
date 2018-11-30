package services

import (
	"github.com/cathalgarvey/fmtless/encoding/json"
	"github.com/oskca/gopherjs-localStorage"
	"honnef.co/go/js/xhr"
)

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

type AuthorizationJson struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
type TokenJson struct {
	Token string `json:"token"`
}

type AuthService struct{}

func (as AuthService) PostLogin(auth AuthorizationJson) (string, *ErrorMsg) {
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
	if IsErrorMsg(data) {
		errMsg := ErrorMsg{}
		err = json.Unmarshal(data, &errMsg)
		if err != nil {
			return "", &ErrorMsg{Error: err.Error()}
		}
		return "", &errMsg
	}
	token := TokenJson{}
	err = json.Unmarshal(data, &token)
	if err != nil {
		return "", &ErrorMsg{Error: err.Error()}
	}
	return token.Token, nil
}

func (as AuthService) GetIsAdmin() (bool, *ErrorMsg) {
	token := as.GetToken()
	req := xhr.NewRequest("GET", "/api/admin/isAdmin")
	req.SetRequestHeader("Content-Type", "application/json")
	req.SetRequestHeader("Authorization", "Bearer "+token)
	err := req.Send(nil)
	if err != nil {
		return false, &ErrorMsg{Error: err.Error()}
	}
	isAdminJson := struct {
		IsAdmin bool `json:"is_admin"`
	}{
		IsAdmin: false,
	}
	err = json.Unmarshal([]byte(req.ResponseText), &isAdminJson)
	if err != nil {
		return false, &ErrorMsg{Error: err.Error()}
	}
	return isAdminJson.IsAdmin, nil
}

func (as AuthService) SaveToken(token string) {
	localStorage.SetItem("token", token)
}

func (as AuthService) RemoveToken() {
	localStorage.SetItem("token", "")
}

func (as AuthService) GetToken() string {
	return localStorage.GetItem("token")
}
