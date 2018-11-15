package admin

import "errors"

var ERR_USERNAME_EMPTY = errors.New("The username is empty.")
var ERR_PASSWORD_NOT_VALID = errors.New("The password is not valid.")
var ERR_TOKEN_PROBLEM = errors.New("Failed to create a token.")
var ERR_URL_NOT_CORRECT = errors.New("The URL is not correct.")
var ERR_MOVIE_NOT_FOUND = errors.New("The movie is not found.")
var ERR_NAME_IS_EMPTY = errors.New("The name is empty.")
var ERR_ITEMS_NOT_ZERO = errors.New("The number of items per page can not be zero.")
var ERR_ITEMS_NOT_LESS_MINUS_ONE = errors.New("The number of items per page can not be less than minus one.")

var ERR_DB_FAILED = func(err error) error {
	return errors.New("The DB failed with this error: " + err.Error())
}
