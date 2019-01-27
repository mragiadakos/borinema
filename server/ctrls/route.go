package ctrls

import (
	"log"
	"os"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/mragiadakos/borinema/server/admin"
	"github.com/mragiadakos/borinema/server/cinema"
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

	err = os.MkdirAll(config.DownloadFolder, 0770)
	if err != nil {
		log.Fatal("Error: " + err.Error())
	}

	cp := utils.NewCommonMoviePlayer()
	aa := admin.NewAdminApi(db, config)
	wa := admin.NewAdminWsApi(config, cp)
	wc := cinema.NewCinemaWsApi(cp)
	cma := cinema.NewCinemaMovieApi(db, config, wa.RequestCurrentTime)

	r.POST("/api/login", aa.Login())
	r.GET("/api/admin/ws", wa.HttpFunc())
	r.GET("/api/player/ws", wc.HttpFunc())

	cinemaGroup := r.Group("/api/cinema")
	cinemaGroup.GET("/movie/:id", cma.GetMovieVideo())
	cinemaGroup.GET("/movie", cma.GetMovieInfo())
	cinemaGroup.POST("/currentTime", cma.RequestCurrentTime())

	adminGroup := r.Group("/api/admin")
	jwtConfig := middleware.JWTConfig{
		Claims:     &utils.JwtCustomClaims{},
		SigningKey: []byte(config.JwtSecret),
	}
	adminGroup.Use(middleware.JWTWithConfig(jwtConfig))
	adminGroup.Use(aa.AuthorizeAdminMiddleware)

	adminGroup.GET("/isAdmin", aa.IsAdmin())
	adminGroup.POST("/movies/link", aa.DownloadMovieLink(wa.SendProgressOfMovie))
	adminGroup.GET("/movies/selected", aa.SelectedMovie())
	adminGroup.DELETE("/movies/selected", aa.RemoveAnySelectedMovie())
	adminGroup.GET("/movies/:id", aa.GetMovie())
	adminGroup.PUT("/movies/:id", aa.UpdateMovie())
	adminGroup.DELETE("/movies/:id", aa.DeleteMovie())
	adminGroup.PUT("/movies/:id/select", aa.SelectMovie())
	adminGroup.GET("/movies", aa.GetMovies())

	r.GET("/", cma.CinemaPage())
	r.GET("/admin", aa.AdminPage())
	r.Static("/admin_panel", "admin_panel")
	r.Static("/cinema_panel", "cinema_panel")
	r.Start(":" + config.Port)
}
