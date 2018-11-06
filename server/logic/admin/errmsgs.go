package admin

import "errors"

var ERR_USERNAME_EMPTY = errors.New("The username is empty.")
var ERR_PASSWORD_NOT_VALID = errors.New("The password is not valid.")
var ERR_TOKEN_PROBLEM = errors.New("Failed to create a token.")
