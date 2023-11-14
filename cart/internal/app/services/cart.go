package services

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"go.opentelemetry.io/otel/trace"

	"route256/cart/internal/app/models"

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

type CartCache interface {
	GetCartItems(ctx context.Context, userId int64) ([]models.CartItemWithInfo, error)
	SetCartItems(ctx context.Context, userId int64, items []models.CartItemWithInfo) error
	Invalidate(ctx context.Context, userId int64) error
}

type cart struct {
	productService ProductService
	lomsService    LomsService
	cartProvider   CartProvider
	cache          CartCache
	tracer         trace.Tracer
}

func NewCartService(
	productService ProductService,
	lomsService LomsService,
	cartProvider CartProvider,
	cache CartCache,
	tracer trace.Tracer,
) Cart {
	return &cart{
		productService: productService,
		lomsService:    lomsService,
		cartProvider:   cartProvider,
		cache:          cache,
		tracer:         tracer,
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
		log.Err(err).Msg("get product error")
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

	err = c.cache.Invalidate(ctx, item.User)
	if err != nil {
		log.Error().Err(err).Int64("user", item.User).Msg("cannot invalidate cache of user cart")
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

	err = c.cache.Invalidate(ctx, item.User)
	if err != nil {
		log.Error().Err(err).Int64("user", item.User).Msg("cannot invalidate cache of user cart")
	}

	return nil
}

func (c *cart) GetItemsByUserId(ctx context.Context, userId int64) ([]models.CartItemWithInfo, uint32, error) {
	if userId == 0 {
		return nil, 0, errors.New("user id is required")
	}

	itemsFromCache, err := c.cache.GetCartItems(ctx, userId)
	if err == nil && itemsFromCache != nil {
		return itemsFromCache, c.getTotalPriceOfItems(itemsFromCache), nil
	}
	if err != nil {
		log.Error().Err(err).Msg("error getting cart items from cache")
	}

	items, err := c.cartProvider.GetItemsByUserId(ctx, userId)
	if err != nil {
		return nil, 0, fmt.Errorf("error fetching users cart: %w", err)
	}

	itemsWithInfo, err := c.getProductsFullInfo(ctx, items)
	if err != nil {
		return nil, 0, err
	}

	err = c.cache.SetCartItems(ctx, userId, itemsWithInfo)
	if err != nil {
		log.Error().Err(err).Msg("error caching cart items")
	}

	return itemsWithInfo, c.getTotalPriceOfItems(itemsWithInfo), nil
}

func (c *cart) getTotalPriceOfItems(itemsWithInfo []models.CartItemWithInfo) uint32 {
	var totalPrice uint32 = 0
	for _, itemWithInfo := range itemsWithInfo {
		totalPrice += itemWithInfo.Price * uint32(itemWithInfo.Count)
	}
	return totalPrice
}

func (c *cart) getProductsFullInfo(ctx context.Context, items []models.CartItem) ([]models.CartItemWithInfo, error) {
	cartItemsWithInfo := make([]models.CartItemWithInfo, len(items))

	// своя реализация golang.org/x/sync/errgroup без ограничения по горутинам
	// cancel необходима для отмены запросов при первой ошибке
	ctx, cancel := context.WithCancel(ctx)
	// wg нужна для того, чтобы дождаться выполнения всех запросов
	wg := &sync.WaitGroup{}
	// нужно для того, чтобы сохранить самую первую ошибку и не перезаписать ее
	errOnce := &sync.Once{}
	// будет хранить ошибку, которая произошла первой при получении инфы о продукте
	var errFetchProduct error

	for i, item := range items {
		// шадоуим переменные, чтобы в горутинах оказались правильные значения
		i, item := i, item
		wg.Add(1)
		// в отдельной горутине делаем запрос в сервис продуктов
		go func() {
			defer wg.Done()
			select {
			// если в какой-либо другой горутине будет ошибка - контекст отменится, нам уже не нужно будет делать запрос
			case <-ctx.Done():
				return
			default:
				itemInfo, err := c.productService.GetProduct(ctx, item.Sku)
				if err != nil {
					log.Err(err).Msg("error fetching product from product service")
					// если произошла ошибка, записываем ее и отменяем контекст, чтобы прекратить выполнение остальных запросов
					// выполняется единожды, чтобы не перезаписать первую ошибку
					errOnce.Do(func() {
						errFetchProduct = err
						cancel()
					})
					return
				}
				// ошибки нет, добавляем итем в результат
				cartItemsWithInfo[i] = models.CartItemWithInfo{
					Sku:   item.Sku,
					Count: item.Count,
					Name:  itemInfo.Name,
					Price: itemInfo.Price,
				}
			}
		}()
	}

	// ждем, пока все горутины завершатся
	wg.Wait()
	// отменяем контекст, чтобы он не протек
	cancel()

	// если мы сохранили ошибку из горутин, возвращаем ее
	if errFetchProduct != nil {
		return nil, errFetchProduct
	}

	return cartItemsWithInfo, nil
}

func (c *cart) Checkout(ctx context.Context, userId int64) (int64, error) {
	ctx, span := c.tracer.Start(ctx, "checkout start", trace.WithSpanKind(trace.SpanKindInternal))
	defer span.End()

	if userId == 0 {
		return 0, errors.New("user id is required")
	}

	_, getItemsSpan := c.tracer.Start(ctx, "get items", trace.WithSpanKind(trace.SpanKindInternal))
	items, err := c.cartProvider.GetItemsByUserId(ctx, userId)
	getItemsSpan.End()
	if err != nil {
		return 0, fmt.Errorf("error fetching users cart: %w", err)
	}
	if len(items) == 0 {
		return 0, ErrCartIsEmpty
	}

	_, createOrderSpan := c.tracer.Start(ctx, "create order", trace.WithSpanKind(trace.SpanKindInternal))
	orderId, err := c.lomsService.CreateOrder(ctx, userId, items)
	createOrderSpan.End()
	if err != nil {
		return 0, fmt.Errorf("error creating order: %w", err)
	}

	_, deleteItemsSpan := c.tracer.Start(ctx, "delete items", trace.WithSpanKind(trace.SpanKindInternal))
	err = c.DeleteItemsByUserId(ctx, userId)
	deleteItemsSpan.End()
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

	err = c.cache.Invalidate(ctx, userId)
	if err != nil {
		log.Error().Err(err).Int64("user", userId).Msg("cannot invalidate cache of user cart")
	}

	return nil
}
