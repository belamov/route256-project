package services

import (
	"context"
	"errors"
	"fmt"

	"github.com/rs/zerolog/log"
	"route256/cart/internal/app/models"
)

type Cart interface {
	AddItem(ctx context.Context, item models.CartItem) error
	DeleteItem(ctx context.Context, item models.CartItem) error
	GetItemsByUserId(ctx context.Context, userId int64) ([]models.CartItemWithInfo, uint32, error)
}

var (
	ErrItemInvalid        = errors.New("item is invalid")
	ErrInsufficientStocks = errors.New("insufficient stocks")
	ErrSkuInvalid         = errors.New("invalid sku")
)

type ProductService interface {
	GetProduct(ctx context.Context, sku uint32) (models.CartItemInfo, error)
}

type StocksService interface {
	GetStocksInfo(ctx context.Context, sku uint32) (uint16, error)
}

type CartStorage interface {
	SaveItem(ctx context.Context, item models.CartItem) error
	DeleteItem(ctx context.Context, item models.CartItem) error
	GetItemsByUserId(ctx context.Context, userId int64) ([]models.CartItem, error)
}

type cartService struct {
	productService ProductService
	stocksService  StocksService
	cartStorage    CartStorage
}

func NewCartService(
	productService ProductService,
	stocksService StocksService,
	cartStorage CartStorage,
) Cart {
	return &cartService{
		productService: productService,
		stocksService:  stocksService,
		cartStorage:    cartStorage,
	}
}

func (c *cartService) AddItem(ctx context.Context, item models.CartItem) error {
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

	stocksAvailable, err := c.stocksService.GetStocksInfo(ctx, item.Sku)
	if err != nil {
		log.Err(err).Msg("get stocks info error")
		return fmt.Errorf("error getting stock for sku: %w", err)
	}
	if stocksAvailable < item.Count {
		return ErrInsufficientStocks
	}

	err = c.cartStorage.SaveItem(ctx, item)
	if err != nil {
		return fmt.Errorf("error adding item to cart: %w", err)
	}

	return nil
}

func (c *cartService) DeleteItem(ctx context.Context, item models.CartItem) error {
	if item.User == 0 {
		return ErrItemInvalid
	}

	err := c.cartStorage.DeleteItem(ctx, item)
	if err != nil {
		return fmt.Errorf("error deleting item from cart: %w", err)
	}

	return nil
}

func (c *cartService) GetItemsByUserId(ctx context.Context, userId int64) ([]models.CartItemWithInfo, uint32, error) {
	if userId == 0 {
		return nil, 0, errors.New("user id is required")
	}

	items, err := c.cartStorage.GetItemsByUserId(ctx, userId)
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
