package auth

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type JwtService struct {
	SecretKey              []byte
	AccessTokenExpiration  time.Duration
	RefreshTokenExpiration time.Duration
}

func NewJwtService(secretKey string, accessTokenExpiration, refreshTokenExpiration time.Duration) *JwtService {
	return &JwtService{
		SecretKey:              []byte(secretKey),
		AccessTokenExpiration:  accessTokenExpiration,
		RefreshTokenExpiration: refreshTokenExpiration,
	}
}

func (s *JwtService) GenerateAccessToken(userDetails *UserDetails) (string, error) {
	return s.generateToken(userDetails, s.AccessTokenExpiration, AccessToken)
}

func (s *JwtService) GenerateRefreshToken(userDetails *UserDetails) (string, error) {
	return s.generateToken(userDetails, s.RefreshTokenExpiration, RefreshToken)
}

func (s *JwtService) ValidateToken(tokenString string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return s.SecretKey, nil
	})
	if err != nil {
		return nil, err
	}
	return token, nil
}

func (s *JwtService) ExtractTokenType(token *jwt.Token) (int, error) {
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return 0, errors.New("invalid token claims")
	}

	tokenType, ok := claims["token_type"].(float64)
	if !ok {
		return 0, errors.New("invalid token type claim")
	}

	return int(tokenType), nil
}
func (s *JwtService) ExtractUserDetails(token *jwt.Token) (*UserDetails, error) {
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("invalid token claims")
	}
	username, ok := claims["sub"].(string)
	if !ok {
		return nil, errors.New("invalid username claim")
	}

	authoritiesInterface, ok := claims["authorities"].([]any)
	if !ok {
		return nil, errors.New("invalid authorities claim")
	}

	authorities := make([]int, len(authoritiesInterface))
	for i, v := range authoritiesInterface {
		floatValue, ok := v.(float64)
		if !ok {
			return nil, errors.New("invalid authority value")
		}
		authorities[i] = int(floatValue)
	}

	return &UserDetails{
		Username:    username,
		Authorities: authorities,
	}, nil
}

func (s *JwtService) IsTokenExpired(token *jwt.Token) (bool, error) {
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return true, errors.New("invalid token claims")
	}

	expiresAt, ok := claims["exp"].(float64)
	if !ok {
		return true, errors.New("invalid expiration claim")
	}

	return time.Now().Unix() > int64(expiresAt), nil
}

func (s *JwtService) generateToken(userDetails *UserDetails, expiration time.Duration, tokenType int) (string, error) {
	claims := jwt.MapClaims{
		"sub":         userDetails.Username,
		"authorities": userDetails.Authorities,
		"token_type":  tokenType,
		"exp":         time.Now().Add(expiration).Unix(),
		"iat":         time.Now().Unix(),
		"jti":         uuid.NewString(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(s.SecretKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
