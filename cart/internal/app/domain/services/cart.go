package services

import (
	"context"
	"errors"
	"fmt"

	"route256/cart/internal/app/domain/models"

	"github.com/rs/zerolog/log"
)

type Cart interface {
	AddItem(ctx context.Context, item models.CartItem) error
	DeleteItem(ctx context.Context, item models.CartItem) error
	GetItemsByUserId(ctx context.Context, userId int64) ([]models.CartItemWithInfo, uint32, error)
	Checkout(ctx context.Context, userId int64) (int64, error)
	DeleteItemsByUserId(ctx context.Context, userId int64) error
}

var (
	ErrItemInvalid        = errors.New("item is invalid")
	ErrInsufficientStocks = errors.New("insufficient stocks")
	ErrSkuInvalid         = errors.New("invalid sku")
	ErrCartIsEmpty        = errors.New("cart is empty")
)

type ProductService interface {
	GetProduct(ctx context.Context, sku uint32) (models.CartItemInfo, error)
}

type LomsService interface {
	GetStocksInfo(ctx context.Context, sku uint32) (uint64, error)
	CreateOrder(ctx context.Context, userId int64, items []models.CartItem) (int64, error)
}

type CartProvider interface {
	SaveItem(ctx context.Context, item models.CartItem) error
	DeleteItem(ctx context.Context, item models.CartItem) error
	GetItemsByUserId(ctx context.Context, userId int64) ([]models.CartItem, error)
	DeleteItemsByUserId(ctx context.Context, userId int64) error
}

type cart struct {
	productService ProductService
	lomsService    LomsService
	cartProvider   CartProvider
}

func NewCartService(
	productService ProductService,
	lomsService LomsService,
	cartProvider CartProvider,
) Cart {
	return &cart{
		productService: productService,
		lomsService:    lomsService,
		cartProvider:   cartProvider,
	}
}

func (c *cart) AddItem(ctx context.Context, item models.CartItem) error {
	if item.Count <= 0 {
		return ErrItemInvalid
	}

	if item.User == 0 {
		return ErrItemInvalid
	}

	_, err := c.productService.GetProduct(ctx, item.Sku)
	if err != nil {
		return ErrSkuInvalid
	}

	stocksAvailable, err := c.lomsService.GetStocksInfo(ctx, item.Sku)
	if err != nil {
		log.Err(err).Msg("get stocks info error")
		return fmt.Errorf("error getting stock for sku: %w", err)
	}
	if stocksAvailable < item.Count {
		return ErrInsufficientStocks
	}

	err = c.cartProvider.SaveItem(ctx, item)
	if err != nil {
		return fmt.Errorf("error adding item to cart: %w", err)
	}

	return nil
}

func (c *cart) DeleteItem(ctx context.Context, item models.CartItem) error {
	if item.User == 0 {
		return ErrItemInvalid
	}

	err := c.cartProvider.DeleteItem(ctx, item)
	if err != nil {
		return fmt.Errorf("error deleting item from cart: %w", err)
	}

	return nil
}

func (c *cart) GetItemsByUserId(ctx context.Context, userId int64) ([]models.CartItemWithInfo, uint32, error) {
	if userId == 0 {
		return nil, 0, errors.New("user id is required")
	}

	items, err := c.cartProvider.GetItemsByUserId(ctx, userId)
	if err != nil {
		return nil, 0, fmt.Errorf("error fetching users cart: %w", err)
	}

	cartItemsWithInfo := make([]models.CartItemWithInfo, 0, len(items))
	var totalPrice uint32 = 0
	for _, item := range items {
		itemInfo, err := c.productService.GetProduct(ctx, item.Sku)
		if err != nil {
			return nil, 0, fmt.Errorf("error fetching product info: %w", err)
		}
		cartItemWithInfo := models.CartItemWithInfo{
			User:  userId,
			Sku:   item.Sku,
			Count: item.Count,
			Name:  itemInfo.Name,
			Price: itemInfo.Price,
		}
		cartItemsWithInfo = append(cartItemsWithInfo, cartItemWithInfo)
		totalPrice += cartItemWithInfo.Price
	}

	return cartItemsWithInfo, totalPrice, nil
}

func (c *cart) Checkout(ctx context.Context, userId int64) (int64, error) {
	if userId == 0 {
		return 0, errors.New("user id is required")
	}

	items, err := c.cartProvider.GetItemsByUserId(ctx, userId)
	if err != nil {
		return 0, fmt.Errorf("error fetching users cart: %w", err)
	}
	if len(items) == 0 {
		return 0, ErrCartIsEmpty
	}

	orderId, err := c.lomsService.CreateOrder(ctx, userId, items)
	if err != nil {
		return 0, fmt.Errorf("error creating order: %w", err)
	}

	err = c.DeleteItemsByUserId(ctx, userId)
	if err != nil {
		return 0, fmt.Errorf("error clearing cart after order creating: %w", err)
	}

	return orderId, nil
}

func (c *cart) DeleteItemsByUserId(ctx context.Context, userId int64) error {
	if userId == 0 {
		return errors.New("user id is required")
	}

	err := c.cartProvider.DeleteItemsByUserId(ctx, userId)
	if err != nil {
		return fmt.Errorf("error clearing cart from storage: %w", err)
	}

	return nil
}
