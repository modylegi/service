package http

import (
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/modylegi/service/pkg/auth"
	mymw "github.com/modylegi/service/pkg/middleware"

	"github.com/rs/zerolog"
	httpSwagger "github.com/swaggo/http-swagger/v2"

	"github.com/modylegi/service/internal/api"
	"github.com/modylegi/service/internal/api/transport/http/handler"
	"github.com/modylegi/service/internal/domain/service"
	"github.com/modylegi/service/pkg/encoding"
)

type apiFunc func(w http.ResponseWriter, r *http.Request) error

func HandleError(h func(w http.ResponseWriter, r *http.Request) error) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := h(w, r); err != nil {
			log := zerolog.Ctx(r.Context())
			var apiErr api.Error
			if errors.As(err, &apiErr) {
				encoding.Encode(w, apiErr.Code, apiErr)
			} else {
				encoding.Encode(w, http.StatusInternalServerError, map[string]interface{}{
					"status":  http.StatusInternalServerError,
					"message": "internal server error",
				})
			}
			log.Err(err).Msg("")
		}
	}
}

type customRoute struct {
	pattern string
	handler apiFunc
}

func registerRoutes(
	cache bool,
	log *zerolog.Logger,
	mux *chi.Mux,
	userSvc service.UserService,
	adminSvc service.AdminService,
	validationSvc service.ValidationService,
	jwtSvc *auth.JwtService,
) {

	mux.Use(mymw.RequestID)
	mux.Use(mymw.RequestLogger(log))
	mux.Use(middleware.Recoverer)

	mux.Handle("/swagger/*", httpSwagger.WrapHandler)

	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	authHandler := handler.NewAuthHandler(userSvc, jwtSvc, log)
	userHandler := handler.NewUserHandler(cache, log, userSvc, validationSvc)
	adminHandler := handler.NewAdminHandler(log, adminSvc, validationSvc)

	publicRoutes := []customRoute{
		{"POST /auth/register", authHandler.Register},
		{"POST /auth/login", authHandler.Login},
		{"POST /auth/refresh", authHandler.RefreshToken},
	}

	userRoutes := []customRoute{
		{"GET /scenario/{user_id}", userHandler.AllBlocksHandler},
		{"GET /scenario/list/{user_id}", userHandler.AllBlocksHandlerIDAndTitle},
		{"GET /scenario/blocks/{user_id}", userHandler.BlockByIDAndOrTitleParam},
		{"GET /scenario/block/{user_id}/list/{block_id}", userHandler.BlockByID},
		{"GET /scenario/block/{user_id}/{block_id}", userHandler.Content},
		{"GET /block/list", adminHandler.BlockIDAndTitleList},
		{"GET /block", adminHandler.BlockByIDAndOrTitle},
		{"GET /block/{block_id}/list", adminHandler.BlockWithoutContentData},
		{"GET /block/{block_id}", adminHandler.Content},
		{"GET /template/list", adminHandler.TemplateList},
		{"GET /template", adminHandler.TemplateByIDNameType},
	}

	adminRoutes := []customRoute{
		{"GET /block/list", adminHandler.BlockIDAndTitleList},
		{"GET /block", adminHandler.BlockByIDAndOrTitle},
		{"GET /block/{block_id}/list", adminHandler.BlockWithoutContentData},
		{"GET /block/{block_id}", adminHandler.Content},
		{"GET /template/list", adminHandler.TemplateList},
		{"GET /template", adminHandler.TemplateByIDNameType},
	}

	for _, route := range publicRoutes {
		mux.HandleFunc(route.pattern, HandleError(route.handler))
	}

	mux.Group(func(r chi.Router) {
		r.Use(mymw.AuthMiddleware(*jwtSvc, log))

		for _, route := range userRoutes {
			r.HandleFunc(route.pattern, HandleError(route.handler))
		}

		r.Group(func(r chi.Router) {
			r.Use(mymw.AdminOnly)
			for _, route := range adminRoutes {
				r.HandleFunc(route.pattern, HandleError(route.handler))
			}
		})

	})
}
