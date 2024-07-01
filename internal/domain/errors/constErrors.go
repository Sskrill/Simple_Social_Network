package domain

import "errors"

var ErrorRefreshTokenExpired = errors.New("refresh token expired")
var ErrorInvalidUsername = errors.New("Invalid user name")
