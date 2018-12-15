package utils

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"github.com/mragiadakos/borinema/server/conf"
)

type JwtCustomClaims struct {
	Admin bool `json:"admin"`
	jwt.StandardClaims
}

func GetTokenAdmin(config conf.Configuration) func() (string, error) {
	return func() (string, error) {
		token := jwt.New(jwt.SigningMethodHS256)
		claims := token.Claims.(jwt.MapClaims)
		claims["admin"] = true
		claims["exp"] = time.Now().Add(time.Hour * 72).Unix()
		return token.SignedString([]byte(config.JwtSecret))
	}
}

func IsAdmin(c echo.Context) bool {
	user, ok := c.Get("user").(*jwt.Token)
	if !ok {
		log.Println("Error: failed to read 'user'", c.Get("Authentication"))
		return false
	}
	claims, ok := user.Claims.(*JwtCustomClaims)
	if !ok {
		log.Println("Error: failed to read 'claims'", user.Claims)
		return false
	}
	admin := claims.Admin

	return admin
}

func ReadAdminTokenString(config conf.Configuration, tokenString string) (bool, *ErrorMsg) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(config.JwtSecret), nil
	})
	errMsg := NewErrorMsg()
	if err != nil {
		errMsg.Error = errors.New("Failed to open the token: " + err.Error())
		log.Println("Error: ", errMsg.Error)
		errMsg.Status = http.StatusUnauthorized
		return false, errMsg
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		errMsg.Error = errors.New("Failed to read the claims from the token.")
		log.Println("Error: ", errMsg.Error)
		errMsg.Status = http.StatusUnauthorized
		return false, errMsg
	}
	isAdminInf, ok := claims["admin"]
	if !ok {
		errMsg.Error = errors.New("Failed to read the key 'admin'")
		log.Println("Error: ", errMsg.Error)
		errMsg.Status = http.StatusUnauthorized
		return false, errMsg
	}
	isAdmin, ok := isAdminInf.(bool)
	if !ok {
		errMsg.Error = errors.New("Failed to read the type from the key 'admin'")
		log.Println("Error: ", errMsg.Error)
		errMsg.Status = http.StatusUnauthorized
		return false, errMsg
	}
	return isAdmin, nil
}
