package domain

import "errors"

var ErrorRefreshTokenExpired = errors.New("refresh token expired")
var ErrorInvalidUsername = errors.New("Invalid user name")

type ErrorResponse struct {
	Message string `json:"ErrorMessage"`
}
