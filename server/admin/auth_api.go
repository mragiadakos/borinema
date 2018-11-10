package admin

import (
	"net/http"

	"github.com/jinzhu/gorm"
	"github.com/labstack/echo"
	"github.com/mragiadakos/borinema/server/conf"
	"github.com/mragiadakos/borinema/server/utils"
)

type adminApi struct {
	db *gorm.DB
}

func NewAdminApi(db *gorm.DB) *adminApi {
	aa := &adminApi{}
	aa.db = db
	return aa
}

func (aa *adminApi) Login(config conf.Configuration) func(echo.Context) error {
	return func(c echo.Context) error {
		opts := AuthorizationAdminInput{}
		c.Bind(&opts)
		isValid := func(opts AuthorizationAdminInput) bool {
			return config.AdminUsername == opts.Username && config.AdminPassword == opts.Password
		}
		al := AdminLogic{}
		sa, errMsg := al.AuthorizeAdmin(opts, isValid, utils.GetTokenAdmin(config))
		if errMsg != nil {
			return c.JSON(errMsg.Status, errMsg.Json())
		}
		return c.JSON(http.StatusAccepted, sa)
	}
}

func (aa *adminApi) IsAdmin() func(echo.Context) error {
	return func(c echo.Context) error {
		al := AdminLogic{}
		output := al.IsAdmin(func() bool {
			return utils.IsAdmin(c)
		})
		return c.JSON(http.StatusOK, output)
	}
}

func (aa *adminApi) AdminPage() func(c echo.Context) error {
	return func(c echo.Context) error {
		return c.File("admin_panel/index.html")
	}
}
