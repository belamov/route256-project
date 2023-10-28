package main

import (
	"context"
	"os"
	"os/signal"
	"route256/loms/internal/app"
	grpcserver "route256/loms/internal/app/grpc/server"
	httpserver "route256/loms/internal/app/http/server"
	"route256/loms/internal/app/services"
	"route256/loms/internal/app/storage/repositories"
	"sync"
	"syscall"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr}).With().Caller().Logger()

	config := app.BuildConfig()

	ctx, _ := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	wg := &sync.WaitGroup{}

	wg.Add(1)
	dbPool, err := repositories.InitPostgresDbConnection(ctx, wg, config)
	if err != nil {
		log.Fatal().Err(err).Msg("Cannot initialize connection to postgres")
		return
	}

	stockPgRepository := repositories.NewStocksPgRepository(dbPool)
	orderPgRepository := repositories.NewOrderPgRepository(dbPool)
	pgTransactor := repositories.NewPgTransactor(dbPool)

	wg.Add(1)
	//eventProducer, err := kafka.NewKafkaEventProducer(
	//	ctx,
	//	wg,
	//	config.KafkaBrokers,
	//	kafka.WithRequiredAcks(sarama.NoResponse),
	//	kafka.WithProducerPartitioner(sarama.NewHashPartitioner),
	//	kafka.WithMaxOpenRequests(1),
	//	kafka.WithMaxRetries(5),
	//	kafka.WithRetryBackoff(10*time.Millisecond),
	//	kafka.WithProducerFlushMessages(3),
	//	kafka.WithProducerFlushFrequency(5*time.Second),
	//)
	//if err != nil {
	//	log.Fatal().Err(err).Msg("Cannot initialize kafka event producer")
	//	return
	//}

	lomsService := services.NewLomsService(
		orderPgRepository,
		stockPgRepository,
		config.AllowedOrderUnpaidTime,
		pgTransactor,
		nil,
	)

	httpServer := httpserver.NewHTTPServer(config.HttpServerAddress, lomsService)
	grpcServer := grpcserver.NewGRPCServer(config.GrpcServerAddress, config.GrpcGatewayServerAddress, lomsService)

	wg.Add(4)

	go httpServer.Run(ctx, wg)
	go lomsService.RunCancelUnpaidOrders(ctx, wg, config.CancelUnpaidOrdersInterval)
	go grpcServer.Run(ctx, wg)
	go grpcServer.RunGateway(ctx, wg)

	wg.Wait()

	log.Info().Msg("goodbye")
}
