package handlers

import (
	"errors"
	"net/http"

	"github.com/modylegi/service/internal/api"
	"github.com/modylegi/service/pkg/encoding"
	"github.com/rs/zerolog"
)

type apiFunc func(w http.ResponseWriter, r *http.Request) error

func Make(h apiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log := zerolog.Ctx(r.Context())
		if err := h(w, r); err != nil {
			var apiErr api.Error
			if errors.As(err, &apiErr) {
				encoding.Encode(w, apiErr.Code, apiErr)
			} else {
				errResp := map[string]interface{}{
					"status":  http.StatusInternalServerError,
					"message": "internal server error",
				}
				encoding.Encode(w, http.StatusInternalServerError, errResp)
			}

			log.Err(err).Msg("")
		}
	}
}
