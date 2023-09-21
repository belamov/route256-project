package server

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"time"

	"route256/loms/internal/app/handlers"
	"route256/loms/internal/app/services"

	"github.com/rs/zerolog/log"
)

type HTTPServer struct {
	server *http.Server
}

func NewHTTPServer(addr string, service services.Loms) *HTTPServer {
	return &HTTPServer{
		server: &http.Server{
			Addr:              addr,
			Handler:           handlers.NewRouter(service),
			ReadHeaderTimeout: 1 * time.Second,
		},
	}
}

func (s *HTTPServer) Run() {
	idleConnectionsClosed := make(chan struct{})
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt)
		<-sigint
		log.Info().Msg("Shutting Down Server")
		if err := s.server.Shutdown(context.Background()); err != nil {
			log.Error().Err(err).Msg("HTTP server Shutdown: ")
		}
		close(idleConnectionsClosed)
	}()

	log.Info().Msgf("Server listening on %s", s.server.Addr)
	err := s.server.ListenAndServe()
	if !errors.Is(err, http.ErrServerClosed) {
		log.Fatal().Err(err).Msg("HTTP server ListenAndServe:")
	}

	<-idleConnectionsClosed
	log.Info().Msg("Goodbye")
}
