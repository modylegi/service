package handler

import (
	"github.com/modylegi/service/pkg/auth"
	"github.com/modylegi/service/pkg/middleware"
	"net/http"

	"github.com/go-playground/validator"
	"github.com/modylegi/service/internal/api"
	"github.com/modylegi/service/internal/domain/service"
	"github.com/modylegi/service/pkg/encoding"
	"github.com/rs/zerolog"
	"golang.org/x/time/rate"
)

type AuthHandler struct {
	userService service.UserService
	jwtService  *auth.JwtService
	log         *zerolog.Logger
	validate    *validator.Validate
	limiter     *rate.Limiter
}

func NewAuthHandler(
	userService service.UserService,
	jwtService *auth.JwtService,
	log *zerolog.Logger,
) *AuthHandler {
	return &AuthHandler{
		userService: userService,
		jwtService:  jwtService,
		log:         log,
		validate:    validator.New(),
		limiter:     rate.NewLimiter(rate.Limit(1), 3),
	}
}

// Register godoc
//	@Summary		Register a new user
//	@Description	Registers a new user in the system
//	@Tags			Authentication
//	@Accept			json
//	@Produce		json
//	@Param			user	body	service.RegisterReq	true	"User registration request"
//	@Success		201		
//	@Failure		400		
//	@Failure		409		
//	@Failure		500		
//	@Router			/auth/register [post]
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	if !h.limiter.Allow() {
		return api.NewError(http.StatusTooManyRequests, api.ErrRateLimitExceeded)
	}

	req, err := encoding.Decode[service.RegisterReq](r)
	if err != nil {
		return api.NewError(http.StatusBadRequest, api.ErrInvalidData)
	}

	if err := h.validate.Struct(req); err != nil {
		return api.NewError(http.StatusBadRequest, api.ErrInvalidData)
	}

	if err := h.userService.Find(ctx, &req); err == nil {
		return api.NewError(http.StatusConflict, api.ErrUserAlreadyExists)
	}

	if err := h.userService.Create(ctx, &req); err != nil {
		h.log.Error().Err(err).Str("username", req.Username).Msg("Failed to create user")
		return api.NewError(http.StatusInternalServerError, api.ErrServer)
	}

	h.log.Info().Str("username", req.Username).Msg("User registered successfully")
	w.WriteHeader(http.StatusCreated)
	return nil
}

// Login godoc
//	@Summary		Login a user
//	@Description	Authenticates a user and returns tokens
//	@Tags			Authentication
//	@Accept			json
//	@Produce		json
//	@Param			user	body		service.LoginReq	true	"User login request"
//	@Success		200		{object}	service.LoginResp
//	@Failure		400		
//	@Failure		401		
//	@Router			/auth/login [post]
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	if !h.limiter.Allow() {
		return api.NewError(http.StatusTooManyRequests, api.ErrRateLimitExceeded)
	}

	req, err := encoding.Decode[service.LoginReq](r)
	if err != nil {
		return api.NewError(http.StatusBadRequest, api.ErrInvalidData)
	}

	if err := h.validate.Struct(req); err != nil {
		return api.NewError(http.StatusBadRequest, api.ErrInvalidData)
	}

	if err := h.userService.Authenticate(ctx, &req); err != nil {
		h.log.Warn().Str("username", req.Username).Msg("Failed login attempt")
		return api.NewError(http.StatusUnauthorized, api.ErrInvalidCredentials)
	}

	userDetails := &auth.UserDetails{
		Username:    req.Username,
		Authorities: []int{auth.UserRole},
	}

	accessToken, err := h.jwtService.GenerateAccessToken(userDetails)
	if err != nil {
		return err
	}

	refreshToken, err := h.jwtService.GenerateRefreshToken(userDetails)
	if err != nil {
		return err
	}

	resp := &service.LoginResp{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}
	h.log.Info().Str("username", req.Username).Msg("User logged in successfully")
	return encoding.Encode(w, http.StatusOK, resp)
}

// RefreshToken godoc
//	@Summary		Refresh access token
//	@Description	Refreshes the access token using the refresh token
//	@Tags			Authentication
//	@Accept			json
//	@Produce		json
//	@Param			Authorization	header		string	true	"Bearer token"
//	@Success		200				{object}	map[string]string
//	@Failure		400				
//	@Failure		401				
//	@Router			/auth/refresh [post]
func (h *AuthHandler) RefreshToken(w http.ResponseWriter, r *http.Request) error {
	if !h.limiter.Allow() {
		return api.NewError(http.StatusTooManyRequests, api.ErrRateLimitExceeded)
	}

	tokenString := middleware.ExtractTokenFromHeader(r)
	if tokenString == "" {
		return api.NewError(http.StatusUnauthorized, auth.ErrUnauthorized)
	}

	token, err := h.jwtService.ValidateToken(tokenString)
	if err != nil {
		return api.NewError(http.StatusUnauthorized, auth.ErrInvalidToken)
	}

	expired, err := h.jwtService.IsTokenExpired(token)
	if err != nil {
		return err
	}

	if expired {
		return api.NewError(http.StatusBadRequest, auth.ErrTokenExpired)
	}

	userDetails, err := h.jwtService.ExtractUserDetails(token)
	if err != nil {

		return api.NewError(http.StatusUnauthorized, auth.ErrInvalidToken)
	}

	tokenType, err := h.jwtService.ExtractTokenType(token)
	if err != nil {
		return err
	}

	if tokenType == auth.AccessToken {
		return api.NewError(http.StatusBadRequest, auth.ErrAccessTokenNotAllowed)
	}

	accessToken, err := h.jwtService.GenerateAccessToken(userDetails)
	if err != nil {
		return err
	}

	h.log.Info().Str("username", userDetails.Username).Msg("Token refreshed successfully")

	return encoding.Encode(w, http.StatusOK, map[string]string{"access_token": accessToken})
}
