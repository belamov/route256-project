package server

import (
	"context"
	"errors"
	"net/http"
	"sync"
	"time"

	"route256/loms/internal/app/services"

	"route256/loms/internal/app/http/handlers"

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

func (s *HTTPServer) Run(ctx context.Context, wg *sync.WaitGroup) {
	go func() {
		<-ctx.Done()
		log.Info().Msg("Shutting Down Http Server")
		if err := s.server.Shutdown(context.Background()); err != nil {
			log.Error().Err(err).Msg("HTTP server Shutdown: ")
		}
		log.Info().Msg("Http Server shut down")
		wg.Done()
	}()

	log.Info().Msgf("Http Server listening on %s", s.server.Addr)
	err := s.server.ListenAndServe()
	if !errors.Is(err, http.ErrServerClosed) {
		log.Fatal().Err(err).Msg("HTTP server ListenAndServe:")
	}
}
