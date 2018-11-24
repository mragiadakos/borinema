package ctrls

import (
	"log"
	"os"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/mragiadakos/borinema/server/admin"
	"github.com/mragiadakos/borinema/server/conf"
	"github.com/mragiadakos/borinema/server/db"
	"github.com/mragiadakos/borinema/server/utils"
)


func Run(config conf.Configuration) {
	r := echo.New()
	db, err := db.NewDB(config.DatabaseFile)
	if err != nil {
		log.Fatal("Error: " + err.Error())
	}
	err = os.MkdirAll(config.DownloadFolder, 0666)
	if err != nil {
		log.Fatal("Error: " + err.Error())
	}

	aa := admin.NewAdminApi(db)
	r.POST("/api/login", aa.Login(config))

	adminGroup := r.Group("/api/admin")
	jwtConfig := middleware.JWTConfig{
		Claims:     &utils.JwtCustomClaims{},
		SigningKey: []byte(config.JwtSecret),
	}
	adminGroup.Use(middleware.JWTWithConfig(jwtConfig))
	adminGroup.Use(aa.AuthorizeAdminMiddleware)

	adminGroup.GET("/isAdmin", aa.IsAdmin())
	adminGroup.POST("/api/admin/movies/link", aa.DownloadMovieLink(config))
	adminGroup.GET("/api/admin/movies/selected", aa.SelectedMovie(config))
	adminGroup.DELETE("/api/admin/movies/selected", aa.RemoveAnySelectedMovie(config))
	adminGroup.GET("/api/admin/movies/:id", aa.GetMovie(config))
	adminGroup.PUT("/api/admin/movies/:id", aa.UpdateMovie(config))
	adminGroup.DELETE("/api/admin/movies/:id", aa.DeleteMovie(config))
	adminGroup.PUT("/api/admin/movies/:id/select", aa.SelectMovie(config))
	adminGroup.PUT("/api/admin/movies", aa.GetMovies(config))

	r.GET("/admin", aa.AdminPage())
	r.Static("/admin_panel", "admin_panel")
	r.Start(":" + config.Port)
}
