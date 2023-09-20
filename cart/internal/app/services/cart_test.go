package services

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
	services "route256/cart/internal/app/mocks"
	"route256/cart/internal/app/models"
)

type CartTestSuite struct {
	suite.Suite
	mockCtrl           *gomock.Controller
	mockCartStorage    *services.MockCartStorage
	mockStocksService  *services.MockStocksService
	mockProductService *services.MockProductService
	cart               Cart
}

type Reporter struct {
	T *testing.T
}

// ensure Reporter implements gomock.TestReporter.
var _ gomock.TestReporter = Reporter{}

// Errorf is equivalent testing.T.Errorf.
func (r Reporter) Errorf(format string, args ...interface{}) {
	r.T.Errorf(format, args...)
}

// Fatalf crashes the program with a panic to allow users to diagnose
// missing expects.
func (r Reporter) Fatalf(format string, args ...interface{}) {
	panic(fmt.Sprintf(format, args...))
}

func (ts *CartTestSuite) SetupSuite() {
	ts.mockCtrl = gomock.NewController(Reporter{ts.T()})
	ts.mockCartStorage = services.NewMockCartStorage(ts.mockCtrl)
	ts.mockStocksService = services.NewMockStocksService(ts.mockCtrl)
	ts.mockProductService = services.NewMockProductService(ts.mockCtrl)
	ts.cart = NewCartService(ts.mockProductService, ts.mockStocksService, ts.mockCartStorage)
}

func TestHandlersTestSuite(t *testing.T) {
	suite.Run(t, new(CartTestSuite))
}

func (ts *CartTestSuite) TestAddItem() {
	ctx := context.Background()
	item := models.CartItem{
		User:  1,
		Sku:   1,
		Count: 1,
	}
	ts.mockStocksService.EXPECT().GetStocksInfo(ctx, item.Sku).Times(1).Return(item.Count+1, nil)
	ts.mockProductService.EXPECT().GetProduct(ctx, item.Sku).Times(1).Return(models.CartItemInfo{}, nil)
	ts.mockCartStorage.EXPECT().SaveItem(ctx, item).Times(1).Return(nil)
	err := ts.cart.AddItem(ctx, item)
	assert.NoError(ts.T(), err)
}

func (ts *CartTestSuite) TestAddItemInvalidCount() {
	ctx := context.Background()
	item := models.CartItem{
		User:  1,
		Sku:   1,
		Count: 0,
	}
	ts.mockStocksService.EXPECT().GetStocksInfo(ctx, item.Sku).Times(1).Return(item.Count+1, nil)
	err := ts.cart.AddItem(ctx, item)
	assert.ErrorIs(ts.T(), err, ErrItemInvalid)
}

func (ts *CartTestSuite) TestAddItemInvalidSku() {
	ctx := context.Background()
	item := models.CartItem{
		User:  1,
		Sku:   1,
		Count: 1,
	}
	ts.mockStocksService.EXPECT().GetStocksInfo(ctx, item.Sku).Times(1).Return(item.Count+1, nil)
	ts.mockProductService.EXPECT().GetProduct(ctx, item.Sku).Times(1).Return(models.CartItemInfo{}, errors.New("not found"))
	err := ts.cart.AddItem(ctx, item)
	assert.ErrorIs(ts.T(), err, ErrSkuInvalid)
}

func (ts *CartTestSuite) TestAddItemInsufficientStock() {
	ctx := context.Background()
	item := models.CartItem{
		User:  1,
		Sku:   1,
		Count: 1,
	}
	ts.mockStocksService.EXPECT().GetStocksInfo(ctx, item.Sku).Times(1).Return(item.Count-1, nil)
	ts.mockProductService.EXPECT().GetProduct(ctx, item.Sku).Times(1).Return(models.CartItemInfo{}, nil)
	err := ts.cart.AddItem(ctx, item)
	assert.ErrorIs(ts.T(), err, ErrInsufficientStocks)
}

func (ts *CartTestSuite) TestDeleteItem() {
	ctx := context.Background()
	item := models.CartItem{
		User:  1,
		Sku:   1,
		Count: 1,
	}

	ts.mockCartStorage.EXPECT().DeleteItem(ctx, item).Return(nil)
	err := ts.cart.DeleteItem(ctx, item)
	assert.NoError(ts.T(), err)
}

func (ts *CartTestSuite) TestGetItemsByUserId() {
	ctx := context.Background()

	items, total, err := ts.cart.GetItemsByUserId(ctx, 0)
	assert.Empty(ts.T(), items)
	assert.Empty(ts.T(), total)
	assert.Error(ts.T(), err)
}

func (ts *CartTestSuite) TestGetItemsByUserIdNoUser() {
	ctx := context.Background()
	var userId int64 = 1
	cartItems := []models.CartItem{
		{
			User:  1,
			Sku:   1,
			Count: 1,
		},
		{
			User:  1,
			Sku:   2,
			Count: 2,
		},
		{
			User:  1,
			Sku:   3,
			Count: 3,
		},
	}
	ts.mockCartStorage.EXPECT().GetItemsByUserId(ctx, userId).Return(cartItems, nil)
	for i, item := range cartItems {
		ts.mockProductService.EXPECT().
			GetProduct(ctx, item.Sku).
			Return(models.CartItemInfo{
				Sku:   item.Sku,
				Name:  strconv.Itoa(i),
				Price: uint32(item.Count),
			}, nil)
	}
	items, total, err := ts.cart.GetItemsByUserId(ctx, userId)
	assert.NoError(ts.T(), err)
	assert.Len(ts.T(), items, len(cartItems))
	assert.Equal(ts.T(), uint32(6), total)
	for i, item := range cartItems {
		assert.Equal(ts.T(), item.Sku, items[i].Sku)
		assert.Equal(ts.T(), item.Count, items[i].Count)
		assert.Equal(ts.T(), strconv.Itoa(i), items[i].Name)
		assert.Equal(ts.T(), uint32(item.Count), items[i].Price)
	}
}
