package ctrls

import (
	"net/http"

	"github.com/labstack/echo"
	"github.com/mragiadakos/borinema/server/conf"
	"github.com/mragiadakos/borinema/server/logic/admin"
)

type AdminApi struct{}

func (aa AdminApi) Login(config conf.Configuration) func(echo.Context) error {
	return func(c echo.Context) error {
		opts := admin.AuthorizationAdminOpts{}
		c.Bind(&opts)
		isValid := func(opts admin.AuthorizationAdminOpts) bool {
			return config.AdminUsername == opts.Username && config.AdminPassword == opts.Password
		}
		sa, errMsg := admin.AuthorizeAdmin(opts, isValid, getTokenAdmin(config))
		if errMsg != nil {
			return c.JSON(errMsg.Status, errMsg.Json())
		}
		return c.JSON(http.StatusAccepted, sa)
	}
}

func (aa AdminApi) IsAdmin() func(echo.Context) error {
	return func(c echo.Context) error {
		output := admin.IsAdmin(func() bool {
			return isAdmin(c)
		})
		return c.JSON(http.StatusOK, output)
	}
}

func (aa AdminApi) AdminPage() func(c echo.Context) error {
	return func(c echo.Context) error {
		return c.File("admin_panel/index.html")
	}
}
