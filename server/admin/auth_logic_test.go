package admin

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAuthorizationAdminFailEmptyUsername(t *testing.T) {
	opts := AuthorizationAdminInput{}
	al := AdminLogic{}
	_, errMsg := al.AuthorizeAdmin(opts, func(opts AuthorizationAdminInput) bool {
		return false
	}, func() (string, error) {
		return "", nil
	})
	assert.NotNil(t, errMsg)
	assert.Equal(t, errMsg.GetVariable("username"), ERR_USERNAME_EMPTY)
}

func TestAuthorizationAdminFailPasswordNotValid(t *testing.T) {
	opts := AuthorizationAdminInput{
		Username: "lallalala",
	}
	al := AdminLogic{}
	_, errMsg := al.AuthorizeAdmin(opts, func(opts AuthorizationAdminInput) bool {
		return false
	}, func() (string, error) {
		return "", nil
	})
	assert.NotNil(t, errMsg)
	assert.Equal(t, errMsg.Error, ERR_PASSWORD_NOT_VALID)
}

func TestAuthorizationAdminSuccess(t *testing.T) {
	opts := AuthorizationAdminInput{
		Username: "admin",
		Password: "admin",
	}
	al := AdminLogic{}

	_, errMsg := al.AuthorizeAdmin(opts, func(opts AuthorizationAdminInput) bool {
		return true
	}, func() (string, error) {
		return "", nil
	})
	assert.Nil(t, errMsg)
}

func TestIsAdminSuccess(t *testing.T) {
	al := AdminLogic{}

	isAdmin := al.IsAdmin(func() bool {
		return true
	})
	assert.True(t, isAdmin.IsAdmin)
}
