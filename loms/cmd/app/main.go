package main

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/IBM/sarama"

	"route256/loms/internal/app"
	"route256/loms/internal/app/events/kafka"
	grpcserver "route256/loms/internal/app/grpc/server"
	httpserver "route256/loms/internal/app/http/server"
	"route256/loms/internal/app/services"
	"route256/loms/internal/app/storage/repositories"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout}).With().Caller().Logger()
	zerolog.SetGlobalLevel(zerolog.ErrorLevel)

	config := app.BuildConfig()

	zerolog.SetGlobalLevel(config.LogLevel)

	ctx, _ := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	wg := &sync.WaitGroup{}

	dbPool, err := repositories.InitPostgresDbConnection(config)
	if err != nil {
		log.Fatal().Err(err).Msg("Cannot initialize connection to postgres")
		return
	}

	stockPgRepository := repositories.NewStocksPgRepository(dbPool)
	orderPgRepository := repositories.NewOrderPgRepository(dbPool)
	outboxRepository := repositories.NewOutboxPgRepository(dbPool)
	pgTransactor := repositories.NewPgTransactor(dbPool)

	wg.Add(1)
	kafkaProducer, err := kafka.NewKafkaEventProducer(
		ctx,
		wg,
		config.KafkaBrokers,
		kafka.WithRequiredAcks(sarama.NoResponse),
		kafka.WithProducerPartitioner(sarama.NewHashPartitioner),
		kafka.WithMaxOpenRequests(1),
		kafka.WithMaxRetries(5),
		kafka.WithRetryBackoff(10*time.Millisecond),
		kafka.WithProducerFlushMessages(3),
		kafka.WithProducerFlushFrequency(5*time.Second),
	)
	if err != nil {
		log.Fatal().Err(err).Msg("Cannot initialize kafka event producer")
		return
	}

	ordersEventProducer := services.NewOutbox(config.OutboxId, kafkaProducer, outboxRepository)

	lomsService := services.NewLomsService(
		orderPgRepository,
		stockPgRepository,
		config.AllowedOrderUnpaidTime,
		pgTransactor,
		ordersEventProducer,
	)

	httpServer := httpserver.NewHTTPServer(config.HttpServerAddress, lomsService)
	grpcServer := grpcserver.NewGRPCServer(config.GrpcServerAddress, config.GrpcGatewayServerAddress, lomsService)

	wg.Add(4)

	go httpServer.Run(ctx, wg)
	go lomsService.RunCancelUnpaidOrders(ctx, wg, config.CancelUnpaidOrdersInterval)
	go grpcServer.Run(ctx, wg)
	go grpcServer.RunGateway(ctx, wg)

	wg.Add(2)
	go ordersEventProducer.StartSendingMessages(ctx, wg, config.OutboxSendInterval)
	go ordersEventProducer.StartRetryingFailedMessages(ctx, wg, config.OutboxRetryInterval)

	wg.Wait()

	log.Info().Msg("Closing pg pool...")
	dbPool.Close()
	log.Info().Msg("pg pool closed")

	log.Info().Msg("goodbye")
}
