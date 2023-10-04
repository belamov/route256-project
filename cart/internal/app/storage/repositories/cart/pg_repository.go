package cart

import (
	"context"
	"fmt"
	"sync"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog/log"
	"route256/cart/internal/app/models"
	"route256/cart/internal/app/services"
	"route256/cart/internal/app/storage/repositories/cart/queries"
)

type pgRepository struct {
	q queries.Querier
}

func NewCartRepository(ctx context.Context, wg *sync.WaitGroup, user string, password string, host string, db string) (services.CartProvider, error) {
	databaseDSN := fmt.Sprintf("postgresql://%s:%s@%s/%s", user, password, host, db)
	dbPool, err := pgxpool.New(ctx, databaseDSN)
	if err != nil {
		return nil, err
	}

	go func() {
		<-ctx.Done()
		log.Info().Msg("Closing cart repository connections...")
		dbPool.Close()
		log.Info().Msg("Cart repository connections closed")
		wg.Done()
	}()

	return &pgRepository{
		q: queries.New(dbPool),
	}, nil
}

func (c pgRepository) SaveItem(ctx context.Context, item models.CartItem) error {
	params := queries.SaveCartItemParams{
		Sku:    int64(item.Sku),
		Count:  int64(item.Count),
		UserID: item.User,
	}

	return c.q.SaveCartItem(ctx, params)
}

func (c pgRepository) DeleteItem(ctx context.Context, item models.CartItem) error {
	params := queries.DeleteCartItemParams{
		UserID: item.User,
		Sku:    int64(item.Sku),
	}

	return c.q.DeleteCartItem(ctx, params)
}

func (c pgRepository) GetItemsByUserId(ctx context.Context, userId int64) ([]models.CartItem, error) {
	itemsFromBd, err := c.q.GetCartItemsByUserId(ctx, userId)
	if err != nil {
		return nil, err
	}

	items := make([]models.CartItem, 0, len(itemsFromBd))
	for _, itemFromBd := range itemsFromBd {
		items = append(items, models.CartItem{
			User:  itemFromBd.UserID,
			Sku:   uint32(itemFromBd.Sku),
			Count: uint64(itemFromBd.Count),
		})
	}

	return items, nil
}

func (c pgRepository) DeleteItemsByUserId(ctx context.Context, userId int64) error {
	return c.q.DeleteCartItemsByUserId(ctx, userId)
}
