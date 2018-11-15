package ctrls

import (
	"log"
	"net/http"
	"os"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/mragiadakos/borinema/server/admin"
	"github.com/mragiadakos/borinema/server/conf"
	"github.com/mragiadakos/borinema/server/db"
	"github.com/mragiadakos/borinema/server/utils"
)

func authorizeAdministrator(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		if utils.IsAdmin(c) {
			return next(c)
		} else {
			return c.JSON(http.StatusUnauthorized, "")
		}
	}
}

func Run(config conf.Configuration) {
	r := echo.New()
	db, err := db.NewDB(config.DatabaseFile)
	if err != nil {
		log.Fatal("Error: " + err.Error())
	}
	err = os.MkdirAll(config.Folder, 0666)
	if err != nil {
		log.Fatal("Error: " + err.Error())
	}

	aa := admin.NewAdminApi(db)
	r.POST("/api/admin/login", aa.Login(config))

	adminGroup := r.Group("/api/admin")
	jwtConfig := middleware.JWTConfig{
		Claims:     &utils.JwtCustomClaims{},
		SigningKey: []byte(config.JwtSecret),
	}
	adminGroup.Use(middleware.JWTWithConfig(jwtConfig))
	adminGroup.Use(authorizeAdministrator)

	adminGroup.GET("/isAdmin", aa.IsAdmin())
	adminGroup.POST("/api/admin/movies/link", aa.DownloadMovieLink(config))
	adminGroup.GET("/api/admin/movies/:id", aa.GetMovie(config))

	r.GET("/admin", aa.AdminPage())
	r.Start(":" + config.Port)
}
