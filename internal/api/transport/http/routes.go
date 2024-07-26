package http

import (
	"errors"
	"net/http"

	"github.com/rs/zerolog"
	httpSwagger "github.com/swaggo/http-swagger/v2"

	"github.com/modylegi/service/internal/api"
	"github.com/modylegi/service/internal/api/transport/http/handlers"
	"github.com/modylegi/service/internal/domain/service"
	"github.com/modylegi/service/pkg/encoding"
)

type apiFunc func(w http.ResponseWriter, r *http.Request) error

func make(h apiFunc) http.HandlerFunc {
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

func registerRoutes(
	cache bool,
	log *zerolog.Logger,
	mux *http.ServeMux,
	userSvc service.UserService,
	adminSvc service.AdminService,
	validationSvc service.ValidationService,
) {
	userHandler := handlers.NewUserHandler(cache, log, userSvc, validationSvc)
	adminHandler := handlers.NewAdminHandler(log, adminSvc, validationSvc)

	mux.Handle("/swagger/", httpSwagger.WrapHandler)

	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	routes := []struct {
		pattern string
		handler apiFunc
	}{
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

	for _, route := range routes {
		mux.HandleFunc(route.pattern, make(route.handler))
	}

}
