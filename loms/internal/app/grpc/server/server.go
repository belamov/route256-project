package server

import (
	"context"
	"errors"
	"net"
	"sync"

	grpcmiddleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpcrecovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"

	lomspb "route256/loms/api/proto"
	"route256/loms/internal/app/services"
)

type GrpcServer struct {
	lomspb.UnimplementedLomsServer
	server        *grpc.Server
	service       services.Loms
	ServerAddress string
}

func NewGRPCServer(
	serverAddress string,
	service services.Loms,
) *GrpcServer {
	s := grpc.NewServer(
		grpc.StreamInterceptor(grpcmiddleware.ChainStreamServer(
			grpcrecovery.StreamServerInterceptor(),
		)),
		grpc.UnaryInterceptor(grpcmiddleware.ChainUnaryServer(
			grpcrecovery.UnaryServerInterceptor(),
		)),
	)
	return &GrpcServer{
		server:        s,
		service:       service,
		ServerAddress: serverAddress,
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

	lomspb.RegisterLomsServer(s.server, s)

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
