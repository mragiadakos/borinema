package ctrls

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/mragiadakos/borinema/server/conf"
	"github.com/mragiadakos/borinema/server/logic/admin"
	"github.com/stretchr/testify/assert"
)

func TestAdminAuthorizationSuccess(t *testing.T) {
	auth := admin.AuthorizationAdminOpts{
		Username: "admin",
		Password: "admin",
	}
	config := conf.Configuration{}
	config.AdminUsername = "admin"
	config.AdminPassword = "admin"
	config.JwtSecret = "secret"

	authJson, _ := json.Marshal(auth)
	e := echo.New()
	jwtConfig := middleware.JWTConfig{
		Claims:     &jwtCustomClaims{},
		SigningKey: []byte(config.JwtSecret),
	}

	req := httptest.NewRequest(http.MethodPost, "/api/admin/login", strings.NewReader(string(authJson)))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	aa := AdminApi{}
	if assert.NoError(t, aa.Login(config)(c)) {
		assert.Equal(t, http.StatusAccepted, rec.Code)
		sa := admin.SuccessAuthorization{}
		err := json.Unmarshal(rec.Body.Bytes(), &sa)
		assert.Nil(t, err)
		assert.NotEqual(t, 0, len(sa.Token))

		req := httptest.NewRequest(http.MethodGet, "/api/admin/isAdmin", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		req.Header.Set(echo.HeaderAuthorization, "Bearer "+sa.Token)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		aa := AdminApi{}
		isAdminHandler := middleware.JWTWithConfig(jwtConfig)(aa.IsAdmin())
		if assert.NoError(t, isAdminHandler(c)) {
			iao := admin.IsAdminOutput{}
			err := json.Unmarshal(rec.Body.Bytes(), &iao)
			assert.Nil(t, err)
			assert.True(t, iao.IsAdmin)
		}
	}
}
