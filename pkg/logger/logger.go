package logger

import (
	"io"
	"os"
	"runtime/debug"
	"time"

	"github.com/rs/zerolog"
)

const (
	envLocal = "local"
	envProd  = "prod"
)

func New(env string) *zerolog.Logger {
	var (
		output   io.Writer
		logLevel zerolog.Level
	)

	switch env {
	case envLocal:
		output = zerolog.ConsoleWriter{
			Out:        os.Stdout,
			TimeFormat: time.RFC3339Nano,
		}
		logLevel = zerolog.DebugLevel
	default:
		output = os.Stdout
		logLevel = zerolog.InfoLevel
	}

	logger := zerolog.New(output).
		Level(logLevel).
		With().
		Timestamp().
		Caller()

	if env == envProd {
		if buildInfo, ok := debug.ReadBuildInfo(); ok {
			logger = logger.
				Int("pid", os.Getpid()).
				Str("go_version", buildInfo.GoVersion)
		}
	}

	log := logger.Logger()
	zerolog.DefaultContextLogger = &log
	return &log
}
