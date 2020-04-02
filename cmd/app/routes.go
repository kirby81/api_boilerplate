package main

import (
	"net/http"

	"github.com/kirby81/api-boilerplate/internal/api"
	"github.com/kirby81/api-boilerplate/internal/auth"

	"github.com/gorilla/mux"
)

func routes(h *auth.Handler) http.Handler {
	router := mux.NewRouter()
	router.Use(api.LogRequest)

	// router.HandleFunc("/", s.handleIndex()).Methods("GET")
	router.HandleFunc("/login", h.HandleLogin()).Methods("POST")
	router.HandleFunc("/signin", h.HandleSignin()).Methods("POST")
	// router.HandleFunc("/protected", s.isAuth(s.handleProtect())).Methods("GET")

	return router
}
