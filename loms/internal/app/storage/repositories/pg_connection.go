package repositories

import (
	"context"
	"fmt"
	"sync"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog/log"
	"route256/loms/internal/app"
)

func InitPostgresDbConnection(ctx context.Context, wg *sync.WaitGroup, config *app.Config) (*pgxpool.Pool, error) {
	databaseDSN := fmt.Sprintf(
		"postgresql://%s:%s@%s/%s",
		config.DbUser,
		config.DbPassword,
		config.DbHost,
		config.DbName,
	)
	dbPool, err := pgxpool.New(ctx, databaseDSN)
	if err != nil {
		return nil, err
	}
	log.Info().Msg("Connected to postgres")

	go func() {
		<-ctx.Done()
		log.Info().Msg("Closing order repository connections...")
		dbPool.Close()
		log.Info().Msg("Order repository connections closed")
		wg.Done()
	}()

	return dbPool, nil
}
