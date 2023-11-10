package services

import (
	"context"

	"route256/cart/internal/app/models"
)

type NoopCache struct{}

func (n NoopCache) GetCartItems(ctx context.Context, userId int64) ([]models.CartItemWithInfo, error) {
	return nil, nil
}

func (n NoopCache) SetCartItems(ctx context.Context, userId int64, items []models.CartItemWithInfo) error {
	return nil
}

func (n NoopCache) Invalidate(ctx context.Context, userId int64) error {
	return nil
}
