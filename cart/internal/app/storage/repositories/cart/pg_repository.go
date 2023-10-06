package cart

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"route256/cart/internal/app/models"
	"route256/cart/internal/app/services"
	"route256/cart/internal/app/storage/repositories/cart/queries"
)

type pgRepository struct {
	dbPool *pgxpool.Pool
}

type txKey struct{}

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

	return c.getQueriesFromContext(ctx).SaveCartItem(ctx, params)
}

func (c pgRepository) DeleteItem(ctx context.Context, item models.CartItem) error {
	params := queries.DeleteCartItemParams{
		UserID: item.User,
		Sku:    int64(item.Sku),
	}

	return c.getQueriesFromContext(ctx).DeleteCartItem(ctx, params)
}

func (c pgRepository) GetItemsByUserId(ctx context.Context, userId int64) ([]models.CartItem, error) {
	itemsFromBd, err := c.getQueriesFromContext(ctx).GetCartItemsByUserId(ctx, userId)
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
	return c.getQueriesFromContext(ctx).DeleteCartItemsByUserId(ctx, userId)
}

func (c pgRepository) getTxFromContext(ctx context.Context) pgx.Tx {
	if tx, ok := ctx.Value(txKey{}).(*pgxpool.Tx); ok {
		return tx
	}

	return nil
}

func (c pgRepository) getQueriesFromContext(ctx context.Context) queries.Querier {
	if tx, ok := ctx.Value(txKey{}).(*pgxpool.Tx); ok {
		return queries.New(tx)
	}

	return queries.New(c.dbPool)
}
