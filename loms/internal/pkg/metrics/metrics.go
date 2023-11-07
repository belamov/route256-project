package metrics

import (
	"context"
	"errors"
	"net/http"
	"sync"

	grpcprom "github.com/grpc-ecosystem/go-grpc-middleware/providers/prometheus"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Metrics struct {
	SrvMetrics   *grpcprom.ServerMetrics
	Reg          *prometheus.Registry
	PanicHandler func(p any) (err error)
}

func InitMetrics() *Metrics {
	srvMetrics := grpcprom.NewServerMetrics(
		grpcprom.WithServerHandlingTimeHistogram(
			grpcprom.WithHistogramBuckets([]float64{0.001, 0.01, 0.1, 0.3, 0.6, 1, 3, 6, 9, 20, 30, 60, 90, 120}),
		),
	)
	reg := prometheus.NewRegistry()
	reg.MustRegister(srvMetrics, collectors.NewGoCollector())

	panicsTotal := promauto.With(reg).NewCounter(prometheus.CounterOpts{
		Name: "grpc_req_panics_recovered_total",
		Help: "Total number of gRPC requests recovered from internal panic.",
	})
	grpcPanicRecoveryHandler := func(p any) (err error) {
		panicsTotal.Inc()
		log.Error().Msgf("recovered from panic: %s", p)
		return status.Errorf(codes.Internal, "%s", p)
	}

	metrics := &Metrics{
		SrvMetrics:   srvMetrics,
		Reg:          reg,
		PanicHandler: grpcPanicRecoveryHandler,
	}

	return metrics
}

func (m *Metrics) RunServer(ctx context.Context, wg *sync.WaitGroup, address string) {
	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.HandlerFor(
		m.Reg,
		promhttp.HandlerOpts{
			// Opt into OpenMetrics e.g. to support exemplars.
			EnableOpenMetrics: true,
		},
	))

	server := &http.Server{
		Addr:    address,
		Handler: mux,
	}

	go func() {
		<-ctx.Done()
		log.Info().Msg("Shutting Down Metrics Server")
		if err := server.Shutdown(context.Background()); err != nil {
			log.Error().Err(err).Msg("Metrics server Shutdown: ")
		}
		log.Info().Msg("Metrics Server shut down")
		wg.Done()
	}()

	log.Info().Msgf("Metrics Server listening on %s", address)
	err := server.ListenAndServe()
	if !errors.Is(err, http.ErrServerClosed) {
		log.Fatal().Err(err).Msg("Metrics server ListenAndServe:")
	}
}
