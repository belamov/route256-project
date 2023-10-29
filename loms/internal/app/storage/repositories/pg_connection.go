package repositories

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog/log"
	"route256/loms/internal/app"
)

func InitPostgresDbConnection(config *app.Config) (*pgxpool.Pool, error) {
	databaseDSN := fmt.Sprintf(
		"postgresql://%s:%s@%s/%s",
		config.DbUser,
		config.DbPassword,
		config.DbHost,
		config.DbName,
	)
	dbPool, err := pgxpool.New(context.Background(), databaseDSN)
	if err != nil {
		return nil, err
	}
	log.Info().Msg("Connected to postgres")

	return dbPool, nil
}
