package handlers

import (
	"net/http"

	"github.com/modylegi/service/pkg/encoding"
)

// HealthHandler godoc
// @Summary		Health check.
// @Description	Health check.
// @Tags			Health
// @Produce		json
// @Success		200
// @Failure		500
// @Router			/health [get]
func HealthHandler(w http.ResponseWriter, r *http.Request) error {
	resp := make(map[string]string)
	resp["message"] = "alive"
	return encoding.Encode(w, http.StatusOK, resp)
}
