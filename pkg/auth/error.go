package auth

import "errors"

var (
	ErrUnauthorized          = errors.New("unauthorized")
	ErrInvalidToken          = errors.New("invalid token")
	ErrTokenExpired          = errors.New("token has expired")
	ErrAccessTokenNotAllowed = errors.New("access token not allowed for refresh")
)
