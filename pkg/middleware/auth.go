package middleware

import (
	"context"
	"errors"
	"github.com/rs/zerolog"
	"net/http"
	"strings"

	"github.com/modylegi/service/pkg/auth"
)

var (
	ErrServer = errors.New("server error")
)

func AuthMiddleware(jwtService auth.JwtService, log *zerolog.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			log.Info().Msg("AuthMiddleware: Starting")

			tokenString := ExtractTokenFromHeader(r)
			if tokenString == "" {
				log.Warn().Msg("AuthMiddleware: No token found in header")
				http.Error(w, auth.ErrUnauthorized.Error(), http.StatusUnauthorized)
				return
			}
			log.Info().Msg("AuthMiddleware: Token extracted, validating")
			token, err := jwtService.ValidateToken(tokenString)
			if err != nil {
				log.Error().Err(err).Msg("AuthMiddleware: Token validation failed")
				http.Error(w, auth.ErrInvalidToken.Error(), http.StatusUnauthorized)
				return
			}

			expired, err := jwtService.IsTokenExpired(token)
			if err != nil {
				http.Error(w, ErrServer.Error(), http.StatusInternalServerError)
				return
			}
			if expired {
				http.Error(w, auth.ErrTokenExpired.Error(), http.StatusUnauthorized)
				return
			}

			log.Info().Msg("AuthMiddleware: Token validated, extracting user details")
			userDetails, err := jwtService.ExtractUserDetails(token)
			if err != nil {
				log.Error().Err(err).Msg("AuthMiddleware: Failed to extract user details")
				http.Error(w, auth.ErrInvalidToken.Error(), http.StatusUnauthorized)
				return
			}

			tokenType, err := jwtService.ExtractTokenType(token)
			if err != nil {
				http.Error(w, ErrServer.Error(), http.StatusInternalServerError)
				return
			}

			if tokenType != auth.AccessToken {
				http.Error(w, auth.ErrInvalidToken.Error(), http.StatusUnauthorized)
				return
			}

			log.Info().Interface("userDetails", userDetails).Msg("AuthMiddleware: User details extracted, adding to context")
			ctx := context.WithValue(r.Context(), "user", userDetails)
			next.ServeHTTP(w, r.WithContext(ctx))
		}
		return http.HandlerFunc(fn)
	}
}

func ExtractTokenFromHeader(r *http.Request) string {
	bearerToken := r.Header.Get("Authorization")
	if len(bearerToken) > 7 && strings.ToUpper(bearerToken[0:7]) == "BEARER " {
		return bearerToken[7:]
	}
	return ""
}
