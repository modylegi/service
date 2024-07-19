package middleware

import (
	"context"
	"net/http"

	"github.com/google/uuid"
)

const (
	requestIDKey    = "requestID"
	requestIDHeader = "X-Request-ID"
)

func (m *Middleware) RequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			requestID := r.Header.Get(requestIDHeader)
			if requestID == "" {
				requestID = uuid.New().String()
			}
			w.Header().Set(requestIDHeader, requestID)
			ctx = context.WithValue(ctx, requestIDKey, requestID)
			r = r.WithContext(ctx)
			next.ServeHTTP(w, r)
		})
}
