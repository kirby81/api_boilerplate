package api

import (
	"net/http"
	"os"
	"strings"

	"github.com/gorilla/handlers"
)

func (s *Server) isAuth(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		auth := r.Header.Get("Authorization")
		if auth == "" {
			s.respond(w, nil, http.StatusBadRequest)
			return
		}

		// Check auth header
		values := strings.Split(auth, " ")
		if len(values) != 2 {
			s.respond(w, nil, http.StatusBadRequest)
			return
		}
		authType, token := values[0], values[1]

		// Check auth type
		if authType != "Bearer" {
			s.respond(w, nil, http.StatusBadRequest)
			return
		}

		// Check token
		if err := s.auth.Parse(token); err != nil {
			s.respond(w, nil, http.StatusUnauthorized)
			return
		}

		h(w, r)
	}
}

func (s *Server) logRequest(h http.Handler) http.Handler {
	return handlers.CombinedLoggingHandler(os.Stdout, h)
}
