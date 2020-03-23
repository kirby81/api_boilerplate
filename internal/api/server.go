package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/mux"

	"github.com/kirby81/api-boilerplate/internal/auth"
	"github.com/kirby81/api-boilerplate/internal/auth/memory"
	"github.com/rs/zerolog/log"
)

type Server struct {
	srv    *http.Server
	config *config
	router *mux.Router
	auth   *auth.Service
}

func NewServer() (*Server, error) {
	config, err := newConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to init config: %w", err)
	}

	auth, err := auth.NewService(memory.NewRepository(), config.TokenSecret)
	if err != nil {
		return nil, fmt.Errorf("failed to init auth service: %w", err)
	}

	return &Server{
		config: config,
		router: mux.NewRouter(),
		auth:   auth,
	}, nil
}

func (s *Server) Run() error {
	s.initSrv()

	done := make(chan bool, 1)
	quit := make(chan os.Signal, 1)

	signal.Notify(quit, os.Interrupt)
	go s.gracefullyShutdown(done, quit)

	log.Info().Msg(fmt.Sprintf("Server is listening at %s", s.srv.Addr))
	if err := s.srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return fmt.Errorf("failed to listen on %s: %w", s.srv.Addr, err)
	}

	<-done
	log.Info().Msg("Server stopped")

	return nil
}

func (s *Server) initSrv() {
	s.routes()
	addr := fmt.Sprintf("%s:%s", s.config.Hostname, s.config.Port)
	s.srv = &http.Server{
		Addr:         addr,
		Handler:      s.router,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  15 * time.Second,
	}
}

func (s *Server) gracefullyShutdown(done chan bool, quit chan os.Signal) {
	<-quit
	log.Info().Msg("Server is shutting down...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	s.srv.SetKeepAlivesEnabled(false)
	if err := s.srv.Shutdown(ctx); err != nil {
		log.Error().Err(err).Msg("Could not gracefully shutdown the server")
		return
	}
	close(done)
}

func (s *Server) respond(w http.ResponseWriter, data interface{}, status int) {
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
