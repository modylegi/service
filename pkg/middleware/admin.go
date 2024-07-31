package middleware

import (
	"net/http"

	"github.com/modylegi/service/pkg/auth"
)

func AdminOnly(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		userDetails := ctx.Value("user").(*auth.UserDetails)
		if !containsAdmin(userDetails.Authorities, auth.AdminRole) {
			http.NotFound(w, r)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func containsAdmin(slice []int, item int) bool {
	for _, element := range slice {
		if element == item {
			return true
		}
	}
	return false
}
