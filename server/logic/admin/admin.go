package admin

import (
	"net/http"
)

type AuthorizationAdminOpts struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type SuccessAuthorization struct {
	Token string `json:"token"`
}

func AuthorizeAdmin(opts AuthorizationAdminOpts,
	isValid func(AuthorizationAdminOpts) bool,
	getToken func() (string, error)) (*SuccessAuthorization, *ErrorMsg) {
	errMsg := NewErrorMsg()
	if len(opts.Username) == 0 {
		errMsg.VariableErrors["username"] = ERR_USERNAME_EMPTY
		errMsg.Status = http.StatusUnprocessableEntity
	}
	if !isValid(opts) {
		errMsg.Error = ERR_PASSWORD_NOT_VALID
		errMsg.Status = http.StatusUnauthorized
	}
	if errMsg.HasErrors() {
		return nil, errMsg
	}
	sa := SuccessAuthorization{}
	var err error
	sa.Token, err = getToken()
	if err != nil {
		errMsg.Error = err
		return nil, errMsg
	}
	return &sa, nil
}

type IsAdminOutput struct {
	IsAdmin bool `json:"is_admin"`
}

func IsAdmin(isAdmin func() bool) IsAdminOutput {
	return IsAdminOutput{IsAdmin: isAdmin()}
}
