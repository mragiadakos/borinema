package admin

import (
	"net/http"

	"github.com/mragiadakos/borinema/server/utils"
)

type AdminLogic struct{}
type AuthorizationAdminInput struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type AuthorizeAdminOutput struct {
	Token string `json:"token"`
}

func (al AdminLogic) AuthorizeAdmin(opts AuthorizationAdminInput,
	isValid func(AuthorizationAdminInput) bool,
	getToken func() (string, error)) (*AuthorizeAdminOutput, *utils.ErrorMsg) {
	errMsg := utils.NewErrorMsg()
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
	sa := AuthorizeAdminOutput{}
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

func (al AdminLogic) IsAdmin(isAdmin func() bool) IsAdminOutput {
	return IsAdminOutput{IsAdmin: isAdmin()}
}
