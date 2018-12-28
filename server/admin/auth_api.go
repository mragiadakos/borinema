package admin

import (
	"net/http"

	"github.com/jinzhu/gorm"
	"github.com/labstack/echo"
	"github.com/mragiadakos/borinema/server/conf"
	"github.com/mragiadakos/borinema/server/utils"
)

type adminApi struct {
	db     *gorm.DB
	config conf.Configuration
}

func NewAdminApi(db *gorm.DB, config conf.Configuration) *adminApi {
	aa := &adminApi{}
	aa.db = db
	aa.config = config
	return aa
}

func (aa *adminApi) AuthorizeAdminMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		if utils.IsAdmin(c) {
			return next(c)
		} else {
			return c.JSON(http.StatusUnauthorized, "")
		}
	}
}

func (aa *adminApi) Login() func(echo.Context) error {
	return func(c echo.Context) error {
		opts := AuthorizationAdminInput{}
		c.Bind(&opts)
		isValid := func(opts AuthorizationAdminInput) bool {
			return aa.config.AdminUsername == opts.Username && aa.config.AdminPassword == opts.Password
		}
		al := AdminLogic{}
		sa, errMsg := al.AuthorizeAdmin(opts, isValid, utils.GetTokenAdmin(aa.config))
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
