package ctrls

import (
	"net/http"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/mragiadakos/borinema/server/conf"
)

func Run(config *conf.Configuration) {
	r := echo.New()
	aa := AdminApi{}

	r.POST("/api/admin/login", aa.Login(*config))

	adminGroup := r.Group("/api/admin")
	jwtConfig := middleware.JWTConfig{
		Claims:     &jwtCustomClaims{},
		SigningKey: []byte(config.JwtSecret),
	}
	adminGroup.Use(middleware.JWTWithConfig(jwtConfig))
	adminGroup.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if isAdmin(c) {
				return next(c)
			} else {
				return c.JSON(http.StatusUnauthorized, "")
			}
		}
	})
	adminGroup.GET("/isAdmin", aa.IsAdmin())

	r.GET("/admin", aa.AdminPage())
	r.Start(":" + config.Port)
}
