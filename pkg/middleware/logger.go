package middleware

import (
	"net/http"
	"time"

	"github.com/rs/zerolog"
)

func RequestLogger(log *zerolog.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			start := time.Now()

			id := GetReqID(ctx)
			var logger zerolog.Logger
			if id != "" {
				logger = log.With().Str(requestIDHeader, id).Logger()
			} else {
				logger = *log
			}

			defer func() {
				logger.
					Info().
					Str("method", r.Method).
					Str("url", r.URL.RequestURI()).
					Dur("elapsed_ms", time.Since(start)).
					Msg("Request")
			}()
			next.ServeHTTP(w, r)
		})
	}
}
