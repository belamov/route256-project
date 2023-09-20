package services

import (
	"context"
	"errors"
	"fmt"
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
	ts.mockProductService.EXPECT().GetProduct(ctx, item.Sku).Times(1).Return(nil)
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
	ts.mockProductService.EXPECT().GetProduct(ctx, item.Sku).Times(1).Return(errors.New("not found"))
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
	ts.mockProductService.EXPECT().GetProduct(ctx, item.Sku).Times(1).Return(nil)
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
