package api

import (
	"encoding/json"
	"net/http"

	"github.com/rs/zerolog/log"
)

func (s *Server) handleIndex() http.HandlerFunc {
	type responseBody struct {
		Msg string `json:"msg"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		s.respond(w, responseBody{Msg: "Index route!"}, http.StatusOK)
	}
}

func (s *Server) handleLogin() http.HandlerFunc {
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
			s.respond(w, nil, http.StatusInternalServerError)
			return
		}

		token, err := s.auth.Login(body.Email, body.Password)
		if err != nil {
			log.Error().Err(err).Msg("Login service error")
			s.respond(w, nil, http.StatusInternalServerError)
			return
		}

		s.respond(w, responseBody{AccessToken: token}, http.StatusOK)
	}
}

func (s *Server) handleSignin() http.HandlerFunc {
	type requestBody struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		var body requestBody

		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			log.Error().Err(err).Msg("Signin decode error")
			s.respond(w, nil, http.StatusInternalServerError)
			return
		}

		err := s.auth.Signin(body.Email, body.Password)
		if err != nil {
			log.Error().Err(err).Msg("Signin service error")
			s.respond(w, nil, http.StatusInternalServerError)
			return
		}
		s.respond(w, nil, http.StatusCreated)
	}
}

func (s *Server) handleProtect() http.HandlerFunc {
	type responseBody struct {
		Msg string `json:"msg"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		s.respond(w, responseBody{Msg: "Protected route!"}, http.StatusOK)
	}
}
