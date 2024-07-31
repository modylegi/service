package middleware

import (
	"context"
	"fmt"
	"net/http"
	"sync/atomic"
)

type ctxKeyRequestID int

const (
	requestIDHeader                 = "X-Request-ID"
	requestIDKey    ctxKeyRequestID = 0
)

var (
	reqID uint64
)

func RequestID(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		requestID := r.Header.Get(requestIDHeader)
		if requestID == "" {
			id := atomic.AddUint64(&reqID, 1)
			requestID = fmt.Sprintf("%d", id)
		}
		ctx = context.WithValue(ctx, requestIDKey, requestID)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
	return http.HandlerFunc(fn)
}

func GetReqID(ctx context.Context) string {
	if ctx == nil {
		return ""
	}
	if reqID, ok := ctx.Value(requestIDKey).(string); ok {
		return reqID
	}
	return ""
}
