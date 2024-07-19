package http

import (
	"fmt"
	"github.com/modylegi/service/internal/domain/service"
	"github.com/modylegi/service/pkg/middleware"
	"github.com/rs/zerolog"
	"net/http"
	"time"
)

func NewServer(
	port string,
	idleTimeout time.Duration,
	readTimeout time.Duration,
	writeTimeout time.Duration,
	cache bool,
	log *zerolog.Logger,
	userSvc service.UserService,
	adminSvc service.AdminService,
	validationSvc service.ValidationService,
) *http.Server {
	mux := http.NewServeMux()
	registerRoutes(cache, log, mux, userSvc, adminSvc, validationSvc)
	var handler http.Handler = mux

	mw := middleware.New(log)
	handler = mw.Logger(handler)

	handler = mw.RequestID(handler)

	return &http.Server{
		Addr:         fmt.Sprintf(":%s", port),
		Handler:      handler,
		IdleTimeout:  idleTimeout,
		ReadTimeout:  readTimeout,
		WriteTimeout: writeTimeout,
	}
}
