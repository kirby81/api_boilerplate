package auth

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/kirby81/api-boilerplate/internal/api"
	"github.com/rs/zerolog/log"
)

type Handler struct {
	auth *Service
}

func NewHandler(auth *Service) (*Handler, error) {
	if auth == nil {
		return nil, errors.New("no auth service provided")
	}

	h := &Handler{
		auth: auth,
	}

	return h, nil
}

func (h *Handler) HandleLogin() http.HandlerFunc {
	type requestBody struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	type responseBody struct {
		AccessToken string `json:"access_token"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		var body requestBody
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			log.Error().Err(err).Msg("Login decode error")
			api.HTTPRespond(w, nil, http.StatusInternalServerError)
			return
		}

		token, err := h.auth.Login(body.Email, body.Password)
		if err != nil {
			log.Error().Err(err).Msg("Login service error")
			api.HTTPRespond(w, nil, http.StatusInternalServerError)
			return
		}

		api.HTTPRespond(w, responseBody{AccessToken: token}, http.StatusOK)
	}
}

func (h *Handler) HandleSignin() http.HandlerFunc {
	type requestBody struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		var body requestBody

		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			log.Error().Err(err).Msg("Signin decode error")
			api.HTTPRespond(w, nil, http.StatusInternalServerError)
			return
		}

		err := h.auth.Signin(body.Email, body.Password)
		if err != nil {
			log.Error().Err(err).Msg("Signin service error")
			api.HTTPRespond(w, nil, http.StatusInternalServerError)
			return
		}
		api.HTTPRespond(w, nil, http.StatusCreated)
	}
}

func (h *Handler) IsAuth(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		auth := r.Header.Get("Authorization")
		if auth == "" {
			api.HTTPRespond(w, nil, http.StatusBadRequest)
			return
		}

		// Check auth header
		values := strings.Split(auth, " ")
		if len(values) != 2 {
			api.HTTPRespond(w, nil, http.StatusBadRequest)
			return
		}
		authType, token := values[0], values[1]

		// Check auth type
		if authType != "Bearer" {
			api.HTTPRespond(w, nil, http.StatusBadRequest)
			return
		}

		// Check token
		if err := h.auth.Parse(token); err != nil {
			api.HTTPRespond(w, nil, http.StatusUnauthorized)
			return
		}

		handler(w, r)
	}
}
