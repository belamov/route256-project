package services

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	services "route256/loms/internal/app/mocks"
	"route256/loms/internal/app/models"

	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
)

type LomsTestSuite struct {
	suite.Suite
	mockCtrl           *gomock.Controller
	mockOrderStorage   *services.MockOrdersStorage
	mockStocksStorage  *services.MockStocksStorage
	mockProductService *services.MockProductService
	loms               Loms
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

func (ts *LomsTestSuite) SetupSuite() {
	ts.mockCtrl = gomock.NewController(Reporter{ts.T()})
	ts.mockStocksStorage = services.NewMockStocksStorage(ts.mockCtrl)
	ts.mockOrderStorage = services.NewMockOrdersStorage(ts.mockCtrl)
	ts.mockProductService = services.NewMockProductService(ts.mockCtrl)
	ts.loms = NewLomsService(ts.mockProductService, ts.mockOrderStorage, ts.mockStocksStorage)
}

func TestHandlersTestSuite(t *testing.T) {
	suite.Run(t, new(LomsTestSuite))
}

func (ts *LomsTestSuite) TestCreateOrder() {
	ctx := context.Background()
	var userId int64 = 1
	orderItems := []models.OrderItem{
		{
			Sku:   1,
			Count: 1,
		},
		{
			Sku:   2,
			Count: 2,
		},
	}

	newOrder := models.Order{
		Id:     1,
		Items:  orderItems,
		Status: models.OrderStatusNew,
	}

	awaitingOrder := models.Order{
		Id:     1,
		Items:  orderItems,
		Status: models.OrderStatusAwaitingPayment,
	}

	ts.mockStocksStorage.EXPECT().Reserve(ctx, gomock.Any()).Return(nil)
	ts.mockOrderStorage.EXPECT().Create(ctx, userId, gomock.Any(), orderItems).Return(newOrder, nil)
	ts.mockOrderStorage.EXPECT().SetStatus(ctx, gomock.Any(), gomock.Any()).Return(awaitingOrder, nil)

	order, err := ts.loms.CreateOrder(ctx, userId, orderItems)
	assert.NoError(ts.T(), err)
	assert.Equal(ts.T(), awaitingOrder, order)
}

func (ts *LomsTestSuite) TestCreateOrderInsufficientStocks() {
	ctx := context.Background()
	var userId int64 = 1
	orderItems := []models.OrderItem{
		{
			Sku:   1,
			Count: 1,
		},
		{
			Sku:   2,
			Count: 2,
		},
	}

	newOrder := models.Order{
		Id:     1,
		Items:  orderItems,
		Status: models.OrderStatusNew,
	}

	failedOrder := models.Order{
		Id:     1,
		Items:  orderItems,
		Status: models.OrderStatusFailed,
	}

	ts.mockOrderStorage.EXPECT().Create(ctx, userId, gomock.Any(), orderItems).Return(newOrder, nil)
	ts.mockStocksStorage.EXPECT().Reserve(ctx, gomock.Any()).Return(ErrInsufficientStocks)
	ts.mockOrderStorage.EXPECT().SetStatus(ctx, gomock.Any(), gomock.Any()).Return(failedOrder, nil)

	order, err := ts.loms.CreateOrder(ctx, userId, orderItems)
	assert.ErrorIs(ts.T(), err, ErrInsufficientStocks)
	assert.Equal(ts.T(), models.Order{}, order)
}
