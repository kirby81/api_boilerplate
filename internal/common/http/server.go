package http

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

type Server struct {
	*http.Server
	ShutdownTimeout time.Duration
}

func NewDefaultServer(addr string, handler http.Handler, logger *log.Logger) *Server {
	return &Server{
		Server: &http.Server{
			Addr:         addr,
			Handler:      handler,
			ErrorLog:     logger,
			ReadTimeout:  5 * time.Second,
			WriteTimeout: 10 * time.Second,
			IdleTimeout:  120 * time.Second,
		},
		ShutdownTimeout: 10 * time.Second,
	}
}

func (s *Server) ListenAndServe() error {
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt, os.Kill)
		<-sigint

		ctx, cancel := context.WithTimeout(context.Background(), s.ShutdownTimeout)
		defer cancel()

		if err := s.Shutdown(ctx); err != nil {
			s.ErrorLog.Printf("HTTP server Shutdown: %v", err)
		}
	}()

	if err := s.Server.ListenAndServe(); err != http.ErrServerClosed {
		return err
	}

	return nil
}
