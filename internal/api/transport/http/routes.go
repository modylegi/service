package http

import (
	"net/http"

	"github.com/rs/zerolog"
	httpSwagger "github.com/swaggo/http-swagger/v2"

	"github.com/modylegi/service/internal/api/transport/http/handlers"
	"github.com/modylegi/service/internal/domain/service"
)

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

	mux.HandleFunc("/swagger/*", httpSwagger.WrapHandler)

	mux.HandleFunc("GET /health", handlers.Make(handlers.HealthHandler))

	mux.HandleFunc("GET /scenario/{user_id}", handlers.Make(userHandler.AllBlocksHandler))
	mux.HandleFunc("GET /scenario/list/{user_id}", handlers.Make(userHandler.AllBlocksHandlerIDAndTitle))
	mux.HandleFunc("GET /scenario/blocks/{user_id}", handlers.Make(userHandler.BlockByIDAndOrTitleParam))
	mux.HandleFunc("GET /scenario/block/{user_id}/list/{block_id}", handlers.Make(userHandler.BlockByID))
	mux.HandleFunc("GET /scenario/block/{user_id}/{block_id}", handlers.Make(userHandler.Content))

	mux.HandleFunc("GET /block/list", handlers.Make(adminHandler.BlockIDAndTitleList))
	mux.HandleFunc("GET /block", handlers.Make(adminHandler.BlockByIDAndOrTitle))
	mux.HandleFunc("GET /block/{block_id}/list", handlers.Make(adminHandler.BlockWithoutContentData))
	mux.HandleFunc("GET /block/{block_id}", handlers.Make(adminHandler.Content))
	mux.HandleFunc("GET /template/list", handlers.Make(adminHandler.TemplateList))
	mux.HandleFunc("GET /template", handlers.Make(adminHandler.TemplateByIDNameType))

}
