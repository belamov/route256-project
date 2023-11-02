package main

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"route256/notifications/internal/app/services"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"route256/notifications/internal/app"
	"route256/notifications/internal/app/events/kafka"
)

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr}).With().Caller().Logger()
	zerolog.SetGlobalLevel(0)

	config := app.BuildConfig()

	zerolog.SetGlobalLevel(0)

	notifier := services.NewNotifier()

	ctx, _ := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	wg := &sync.WaitGroup{}

	kafkaConfig := kafka.Config{
		Brokers:         config.KafkaBrokers,
		TopicNames:      config.TopicNames,
		ConsumerGroupId: config.ConsumerGroupId,
	}

	wg.Add(1)
	kafkaEventListener, err := kafka.NewOrderStatusChangedNotifier(ctx, wg, kafkaConfig, notifier)
	if err != nil {
		panic(err)
	}

	wg.Add(1)
	kafkaEventListener.Run(ctx, wg)

	<-ctx.Done()
	wg.Wait()
	log.Info().Msg("goodbye")
}
