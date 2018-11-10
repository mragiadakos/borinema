package admin

import "errors"

var ERR_USERNAME_EMPTY = errors.New("The username is empty.")
var ERR_PASSWORD_NOT_VALID = errors.New("The password is not valid.")
var ERR_TOKEN_PROBLEM = errors.New("Failed to create a token.")
var ERR_URL_NOT_CORRECT = errors.New("The URL is not correct.")
var ERR_MOVIE_NOT_FOUND = errors.New("The movie is not found.")
var ERR_DB_FAILED = func(err error) error {
	return errors.New("The DB failed with this error: " + err.Error())
}
