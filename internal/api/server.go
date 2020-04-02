package api

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/rs/zerolog/log"
)

type Server struct {
	srv    *http.Server
	router http.Handler
}

func NewServer(router http.Handler) (*Server, error) {
	if router == nil {
		return nil, errors.New("nothing to route")
	}

	return &Server{
		router: router,
	}, nil
}

func (s *Server) Run(host, port string) error {
	s.initSrv(host, port)

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

func (s *Server) initSrv(host, port string) {
	addr := fmt.Sprintf("%s:%s", host, port)
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
