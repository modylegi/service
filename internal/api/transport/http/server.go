package http

import (
	"fmt"
	"net/http"
	"time"

	"github.com/modylegi/service/internal/domain/service"
	"github.com/modylegi/service/pkg/middleware"
	"github.com/rs/zerolog"
)

type Server struct {
	*http.Server
	middlewares []func(http.Handler) http.Handler
}

func NewServer(
	port int,
	idleTimeout, readTimeout, writeTimeout time.Duration,
	cache bool,
	log *zerolog.Logger,
	userSvc service.UserService,
	adminSvc service.AdminService,
	validationSvc service.ValidationService,
) *Server {
	mux := http.NewServeMux()
	registerRoutes(cache, log, mux, userSvc, adminSvc, validationSvc)

	server := &Server{
		Server: &http.Server{
			Addr:         fmt.Sprintf(":%d", port),
			Handler:      mux,
			IdleTimeout:  idleTimeout,
			ReadTimeout:  readTimeout,
			WriteTimeout: writeTimeout,
		},
	}

	mw := middleware.New(log)
	server.Use(mw.Logger)
	server.Use(mw.RequestID)

	return server
}

func (s *Server) Use(middleware func(http.Handler) http.Handler) {
	s.middlewares = append(s.middlewares, middleware)
}

func (s *Server) Start() error {
	var handler http.Handler = s.Handler
	for i := len(s.middlewares) - 1; i >= 0; i-- {
		handler = s.middlewares[i](handler)
	}
	s.Handler = handler

	return s.ListenAndServe()
}
