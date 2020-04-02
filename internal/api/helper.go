package api

import (
	"encoding/json"
	"net/http"

	"github.com/rs/zerolog/log"
)

func HTTPRespond(w http.ResponseWriter, data interface{}, status int) {
	w.WriteHeader(status)
	if data != nil {
		w.Header().Add("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(data); err != nil {
			log.Error().Err(err).Msg("Failed to encode response data")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
}
