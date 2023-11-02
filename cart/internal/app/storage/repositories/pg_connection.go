package repositories

import (
	"context"
	"fmt"
	"github.com/exaring/otelpgx"
	"sync"

	"route256/cart/internal/app"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog/log"
)

func InitPostgresDbConnection(ctx context.Context, wg *sync.WaitGroup, config *app.Config) (*pgxpool.Pool, error) {
	databaseDSN := fmt.Sprintf(
		"postgresql://%s:%s@%s/%s",
		config.DbUser,
		config.DbPassword,
		config.DbHost,
		config.DbName,
	)
	cfg, err := pgxpool.ParseConfig(databaseDSN)
	if err != nil {
		return nil, fmt.Errorf("parse config: %w", err)
	}

	cfg.ConnConfig.Tracer = otelpgx.NewTracer()

	dbPool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		return nil, fmt.Errorf("create connection pool: %w", err)
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
