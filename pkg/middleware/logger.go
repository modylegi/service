package middleware

import (
	"net/http"
	"time"

	"github.com/rs/zerolog"
)

func (m *Middleware) Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			requestID, ok := r.Context().Value(requestIDKey).(string)
			if !ok {
				requestID = "unknown"
			}

			var log zerolog.Logger
			if requestID != "unknown" {
				log = m.log.With().Str(requestIDHeader, requestID).Logger()
			} else {
				log = *m.log
			}

			ctx := log.WithContext(r.Context())
			r = r.WithContext(ctx)

			defer func() {
				zerolog.Ctx(ctx).
					Info().
					Str("method", r.Method).
					Str("url", r.URL.RequestURI()).
					Dur("elapsed_ms", time.Since(start)).
					Msg("")
			}()

			next.ServeHTTP(w, r)
		})
}
