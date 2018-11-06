package ctrls

import (
	"fmt"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"github.com/mragiadakos/borinema/server/conf"
)

type jwtCustomClaims struct {
	Admin bool `json:"admin"`
	jwt.StandardClaims
}

func getTokenAdmin(config conf.Configuration) func() (string, error) {
	return func() (string, error) {
		token := jwt.New(jwt.SigningMethodHS256)
		claims := token.Claims.(jwt.MapClaims)
		claims["admin"] = true
		claims["exp"] = time.Now().Add(time.Hour * 72).Unix()
		return token.SignedString([]byte(config.JwtSecret))
	}
}
func isAdmin(c echo.Context) bool {
	user, ok := c.Get("user").(*jwt.Token)
	if !ok {
		fmt.Println("user", c.Get("Authentication"))
		return false
	}
	claims, ok := user.Claims.(*jwtCustomClaims)
	if !ok {
		fmt.Println("claims", user.Claims)
		return false
	}
	admin := claims.Admin

	return admin
}
