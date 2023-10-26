package main

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"route256/notifications/internal/app"
	"route256/notifications/internal/app/events/kafka"
)

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr}).With().Caller().Logger()

	config := app.BuildConfig()

	ctx, _ := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	wg := &sync.WaitGroup{}

	wg.Add(1)
	kafkaEventListener, err := kafka.NewConsumer(ctx, wg, config.KafkaBrokers, config.TopicNames, config.ConsumerGroupId)
	if err != nil {
		panic(err)
	}

	wg.Add(1)
	kafkaEventListener.StartConsuming(ctx, wg)

	<-ctx.Done()
	wg.Wait()
	log.Info().Msg("goodbye")
}
