package admin

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAuthorizationAdminFailEmptyUsername(t *testing.T) {
	opts := AuthorizationAdminOpts{}
	_, errMsg := AuthorizeAdmin(opts, func(opts AuthorizationAdminOpts) bool {
		return false
	}, func() (string, error) {
		return "", nil
	})
	assert.NotNil(t, errMsg)
	assert.Equal(t, errMsg.GetVariable("username"), ERR_USERNAME_EMPTY)
}

func TestAuthorizationAdminFailPasswordNotValid(t *testing.T) {
	opts := AuthorizationAdminOpts{
		Username: "lallalala",
	}
	_, errMsg := AuthorizeAdmin(opts, func(opts AuthorizationAdminOpts) bool {
		return false
	}, func() (string, error) {
		return "", nil
	})
	assert.NotNil(t, errMsg)
	assert.Equal(t, errMsg.Error, ERR_PASSWORD_NOT_VALID)
}

func TestAuthorizationAdminSuccess(t *testing.T) {
	opts := AuthorizationAdminOpts{
		Username: "admin",
		Password: "admin",
	}
	_, errMsg := AuthorizeAdmin(opts, func(opts AuthorizationAdminOpts) bool {
		return true
	}, func() (string, error) {
		return "", nil
	})
	assert.Nil(t, errMsg)
}

func TestIsAdminSuccess(t *testing.T) {
	isAdmin := IsAdmin(func() bool {
		return true
	})
	assert.True(t, isAdmin.IsAdmin)
}
