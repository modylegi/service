package http

import (
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/modylegi/service/pkg/auth"

	"github.com/modylegi/service/internal/domain/service"
	"github.com/rs/zerolog"
)

type Server struct {
	*http.Server
}

func NewServer(
	port int,
	idleTimeout, readTimeout, writeTimeout time.Duration,
	cache bool,
	log *zerolog.Logger,
	userSvc service.UserService,
	adminSvc service.AdminService,
	validationSvc service.ValidationService,
	jwtSvc *auth.JwtService,
) *Server {
	mux := chi.NewMux()
	registerRoutes(cache, log, mux, userSvc, adminSvc, validationSvc, jwtSvc)

	return &Server{
		Server: &http.Server{
			Addr:         fmt.Sprintf(":%d", port),
			Handler:      mux,
			IdleTimeout:  idleTimeout,
			ReadTimeout:  readTimeout,
			WriteTimeout: writeTimeout,
		},
	}
}

func (s *Server) Start() error {
	var handler http.Handler = s.Handler
	s.Handler = handler
	return s.ListenAndServe()
}
