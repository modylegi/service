package app

import (
	"context"
	"fmt"
	"github.com/modylegi/service/pkg/rdclient"
	"io"
	"os"
	"os/signal"
	"time"

	"github.com/modylegi/service/internal/api/transport/http"
	"github.com/modylegi/service/internal/service"
	"github.com/modylegi/service/pkg/logger"
	"github.com/modylegi/service/pkg/pgconn"
)

func Run(ctx context.Context, w io.Writer, args []string) error {
	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt)
	defer cancel()

	cfg := getConfig()

	log := logger.New(cfg.Env)

	// postgres setup
	pgCfg := pgconn.NewConfig().
		WithDatabase(cfg.DBDatabase).
		WithPassword(cfg.DBPassword).
		WithUsername(cfg.DBUsername).
		WithHost(cfg.DBHost).
		WithPort(cfg.DBPort)
	db, err := pgconn.NewConnection(ctx, pgCfg)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			fmt.Fprintf(os.Stderr, "Error closing database connection: %v\n", err)
		}
	}()

	// redis setup
	rdCfg := rdclient.NewConfig().
		WithHost(cfg.CacheHost).
		WithPort(cfg.CachePort).
		WithPassword(cfg.CachePassword)
	rd, err := rdclient.NewClient(ctx, rdCfg)
	if err != nil {
		return fmt.Errorf("failed to connect to redis: %w", err)
	}
	defer func() {
		if err := rd.Close(); err != nil {
			fmt.Fprintf(os.Stderr, "Error closing redis client: %v\n", err)
		}
	}()

	userSvc := service.NewUserService(db, rd)
	adminSvc := service.NewAdminService(db)
	validationSvc := service.NewValidationService(db)

	server := http.NewServer(
		cfg.ServerPort,
		cfg.ServerIdleTimeout,
		cfg.ServerReadTimeout,
		cfg.ServerWriteTimeout,
		cfg.Cache,
		log,
		userSvc,
		adminSvc,
		validationSvc,
	)
	log.Info().Msgf("Listening on %s", server.Addr)

	ch := make(chan error, 1)

	go func() {
		if err := server.ListenAndServe(); err != nil {
			ch <- fmt.Errorf("error listening and serving: %w", err)
		}
		close(ch)
	}()

	select {
	case err = <-ch:
		return err
	case <-ctx.Done():
		timeout, cancel := context.WithTimeout(context.Background(), time.Second*20)
		defer cancel()
		return server.Shutdown(timeout)
	}
}
