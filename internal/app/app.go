package app

import (
	"context"
	"fmt"
	"github.com/modylegi/service/pkg/auth"
	"io"
	"os"
	"os/signal"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/modylegi/service/internal/config"
	"github.com/modylegi/service/pkg/rdclient"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog"

	"github.com/modylegi/service/internal/api/transport/http"
	"github.com/modylegi/service/internal/service"
	"github.com/modylegi/service/pkg/logger"
	"github.com/modylegi/service/pkg/pgconn"
)

func Run(ctx context.Context, w io.Writer, args []string) error {
	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt)
	defer cancel()

	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("config load failed: %w", err)
	}

	log := logger.New(cfg.Env)

	db, err := setupDatabase(ctx, cfg)
	if err != nil {
		return fmt.Errorf("database setup failed: %w", err)
	}
	defer closeResource(db, "database", log)

	var rd *redis.Client
	if cfg.UseCache {
		rd, err := setupRedis(ctx, cfg)
		if err != nil {
			return fmt.Errorf("redis setup failed: %w", err)
		}
		defer closeResource(rd, "redis", log)
	}

	server := setupServer(cfg, log, db, rd)

	return runServer(ctx, server, log)
}

func setupDatabase(ctx context.Context, cfg *config.Config) (*sqlx.DB, error) {
	pgCfg := pgconn.NewConfig().
		WithHost(cfg.DB.Host).
		WithPort(cfg.DB.Port).
		WithDatabase(cfg.DB.Database).
		WithUsername(cfg.DB.Username).
		WithPassword(cfg.DB.Password)

	return pgconn.NewConnection(ctx, pgCfg)
}

func setupRedis(ctx context.Context, cfg *config.Config) (*redis.Client, error) {
	rdCfg := rdclient.NewConfig().
		WithHost(cfg.Cache.Host).
		WithPort(cfg.Cache.Port).
		WithPassword(cfg.Cache.Password)

	return rdclient.NewClient(ctx, rdCfg)
}

func setupServer(cfg *config.Config, log *zerolog.Logger, db *sqlx.DB, rd *redis.Client) *http.Server {
	userSvc := service.NewUserService(db, rd)
	adminSvc := service.NewAdminService(db)
	validationSvc := service.NewValidationService(db)
	jwtSvc := auth.NewJwtService("", time.Minute*15, time.Hour*24*7)

	return http.NewServer(
		cfg.HttpServer.Port,
		cfg.HttpServer.IdleTimeout,
		cfg.HttpServer.ReadTimeout,
		cfg.HttpServer.WriteTimeout,
		cfg.UseCache,
		log,
		userSvc,
		adminSvc,
		validationSvc,
		jwtSvc,
	)
}

func runServer(ctx context.Context, server *http.Server, log *zerolog.Logger) error {
	log.Info().Msgf("Listening on %s", server.Addr)

	errCh := make(chan error, 1)
	go func() {
		errCh <- server.Start()
	}()

	select {
	case err := <-errCh:
		return fmt.Errorf("server error: %w", err)
	case <-ctx.Done():
		log.Info().Msg("Shutdown signal received")
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
		defer cancel()
		return server.Shutdown(shutdownCtx)
	}
}

func closeResource(closer io.Closer, resourceName string, log *zerolog.Logger) {
	if err := closer.Close(); err != nil {
		log.Error().Err(err).Msgf("Error closing %s connection", resourceName)
	}
}
