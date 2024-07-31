package api

import (
	"errors"
	"fmt"
)

var (
	ErrNoUserID                         = errors.New("параметр user_id не указан")
	ErrNoBlockID                        = errors.New("параметр block_id не указан")
	ErrParamsDoNotMatchEachOther        = errors.New("параметры друг другу не соответствуют")
	ErrNoScenario                       = errors.New("пользователь не привязан ни к какому сценарию")
	ErrBlockIDNotMatchesUserScenario    = errors.New("block_id не соответствует сценарию пользователя")
	ErrBlockTitleNotMatchesUserScenario = errors.New("name не соответствует сценарию пользователя")
	ErrNoBlockIDAndTitle                = errors.New("не указан ни один из параметров block_id, name")
	ErrNoContentIDAndNameAndType        = errors.New("не указан ни один из параметров content_id, name, content_type")
	ErrNoTemplateIDAndNameAndType       = errors.New("не указан ни один из параметров template_id, name, content_type")
	ErrInvalidData                      = errors.New("invalid data")
	ErrUnauthorized                     = errors.New("unauthorized")
	ErrInvalidCredentials               = errors.New("invalid credentials")
	ErrInvalidToken                     = errors.New("invalid token")
	ErrTokenExpired                     = errors.New("token has expired")
	ErrAccessTokenNotAllowed            = errors.New("access token not allowed for refresh")
	ErrUserAlreadyExists                = errors.New("user already exists ")
	ErrServer                           = errors.New("server error")
	ErrRateLimitExceeded                = errors.New("rate limit exceeded")
	ErrWeakPassword                     = errors.New("password too weak")
	ErrTokenRevoked                     = errors.New("token has been revoked")
)

type Error struct {
	Code    int `json:"status_code"`
	Message any `json:"message"`
}

func (e Error) Error() string {
	return fmt.Sprintf("code: %d, message: %v", e.Code, e.Message)
}

func NewError(code int, err error) Error {
	return Error{
		Code:    code,
		Message: err.Error(),
	}
}
