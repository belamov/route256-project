package server

import (
	"context"
	"errors"
	"net"
	"net/http"
	"sync"

	"route256/loms/internal/pkg/metrics"

	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc/credentials/insecure"

	"route256/loms/internal/app/grpc/pb"
	"route256/loms/internal/app/services"

	grpcmiddleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpcrecovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
)

type GrpcServer struct {
	pb.UnimplementedLomsServer
	server               *grpc.Server
	service              services.Loms
	ServerAddress        string
	GatewayServerAddress string
}

func NewGRPCServer(
	serverAddress string,
	gatewayServerAddress string,
	service services.Loms,
	metrics *metrics.Metrics,
) *GrpcServer {
	s := grpc.NewServer(
		grpc.StreamInterceptor(grpcmiddleware.ChainStreamServer(
			grpcrecovery.StreamServerInterceptor(),
		)),
		grpc.UnaryInterceptor(grpcmiddleware.ChainUnaryServer(
			grpcrecovery.UnaryServerInterceptor(grpcrecovery.WithRecoveryHandler(metrics.PanicHandler)),
			metrics.SrvMetrics.UnaryServerInterceptor(),
		)),
		grpc.StatsHandler(otelgrpc.NewServerHandler()),
	)
	return &GrpcServer{
		server:               s,
		service:              service,
		ServerAddress:        serverAddress,
		GatewayServerAddress: gatewayServerAddress,
	}
}

func (s *GrpcServer) Run(ctx context.Context, wg *sync.WaitGroup) {
	go func() {
		<-ctx.Done()
		log.Info().Msg("Stopping grpc server")
		s.server.GracefulStop()
		log.Info().Msg("Stopped grpc server")
		wg.Done()
	}()

	pb.RegisterLomsServer(s.server, s)

	listen, err := net.Listen("tcp", s.ServerAddress)
	if err != nil {
		log.Fatal().Err(err)
		return
	}

	log.Info().Msgf("Grpc Server listening on %s", s.ServerAddress)

	err = s.server.Serve(listen)
	if err != nil && !errors.Is(err, grpc.ErrServerStopped) {
		log.Error().Err(err).Msg("Grpc Server fail")
	}
}

func (s *GrpcServer) RunGateway(ctx context.Context, wg *sync.WaitGroup) {
	conn, err := grpc.DialContext(
		ctx,
		s.ServerAddress,
		grpc.WithBlock(),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to dial server")
		return
	}

	gwmux := runtime.NewServeMux()
	err = pb.RegisterLomsHandler(ctx, gwmux, conn)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to register gateway")
		return
	}

	gwServer := &http.Server{Addr: s.GatewayServerAddress, Handler: gwmux}

	go func() {
		<-ctx.Done()
		log.Info().Msg("shutting down grpc gateway server")
		if err := conn.Close(); err != nil {
			log.Error().Err(err).Msg("grpc gateway conn close: ")
		}
		if err := gwServer.Shutdown(context.Background()); err != nil {
			log.Error().Err(err).Msg("grpc gateway server shutdown: ")
		}
		log.Info().Msg("grpc gateway server shut down")
		wg.Done()
	}()

	log.Info().Msgf("grpc gateway Server listening on %s", s.GatewayServerAddress)
	err = gwServer.ListenAndServe()
	if !errors.Is(err, http.ErrServerClosed) {
		log.Fatal().Err(err).Msg("grpc gateway server ListenAndServe:")
	}
}
