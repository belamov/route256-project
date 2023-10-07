package repositories

import (
	"context"

	"route256/cart/internal/app/storage/repositories/queries"

	"github.com/jackc/pgx/v5/pgxpool"
	"route256/cart/internal/app/models"
	"route256/cart/internal/app/services"
)

type pgRepository struct {
	dbPool *pgxpool.Pool
}

func NewCartRepository(dbPool *pgxpool.Pool) services.CartProvider {
	return &pgRepository{
		dbPool: dbPool,
	}
}

func (c pgRepository) SaveItem(ctx context.Context, item models.CartItem) error {
	params := queries.SaveCartItemParams{
		Sku:    int64(item.Sku),
		Count:  int64(item.Count),
		UserID: item.User,
	}

	return queries.New(c.dbPool).SaveCartItem(ctx, params)
}

func (c pgRepository) DeleteItem(ctx context.Context, item models.CartItem) error {
	params := queries.DeleteCartItemParams{
		UserID: item.User,
		Sku:    int64(item.Sku),
	}

	return queries.New(c.dbPool).DeleteCartItem(ctx, params)
}

func (c pgRepository) GetItemsByUserId(ctx context.Context, userId int64) ([]models.CartItem, error) {
	itemsFromBd, err := queries.New(c.dbPool).GetCartItemsByUserId(ctx, userId)
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
	return queries.New(c.dbPool).DeleteCartItemsByUserId(ctx, userId)
}
