package services

import (
	"context"
	"fmt"
	"testing"

	"route256/loms/internal/app/domain/models"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
)

type LomsTestSuite struct {
	suite.Suite
	mockCtrl           *gomock.Controller
	mockOrdersProvider *MockOrdersProvider
	mockStocksProvider *MockStocksProvider
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
	ts.mockStocksProvider = NewMockStocksProvider(ts.mockCtrl)
	ts.mockOrdersProvider = NewMockOrdersProvider(ts.mockCtrl)
	ts.loms = NewLomsService(ts.mockOrdersProvider, ts.mockStocksProvider)
}

func TestLomsTestSuite(t *testing.T) {
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

	ts.mockStocksProvider.EXPECT().Reserve(ctx, gomock.Any()).Return(nil)
	ts.mockOrdersProvider.EXPECT().Create(ctx, userId, gomock.Any(), orderItems).Return(newOrder, nil)
	ts.mockOrdersProvider.EXPECT().SetStatus(ctx, gomock.Any(), gomock.Any()).Return(awaitingOrder, nil)

	order, err := ts.loms.OrderCreate(ctx, userId, orderItems)
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

	ts.mockOrdersProvider.EXPECT().Create(ctx, userId, gomock.Any(), orderItems).Return(newOrder, nil)
	ts.mockStocksProvider.EXPECT().Reserve(ctx, gomock.Any()).Return(ErrInsufficientStocks)
	ts.mockOrdersProvider.EXPECT().SetStatus(ctx, gomock.Any(), gomock.Any()).Return(failedOrder, nil)

	order, err := ts.loms.OrderCreate(ctx, userId, orderItems)
	assert.ErrorIs(ts.T(), err, ErrInsufficientStocks)
	assert.Equal(ts.T(), models.Order{}, order)
}

func (ts *LomsTestSuite) TestGetOrderById() {
	ctx := context.Background()
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

	foundOrder := models.Order{
		Id:     1,
		Items:  orderItems,
		Status: models.OrderStatusNew,
	}

	ts.mockOrdersProvider.EXPECT().GetOrderByOrderId(ctx, foundOrder.Id).Return(foundOrder, nil)

	order, err := ts.loms.OrderInfo(ctx, foundOrder.Id)
	assert.NoError(ts.T(), err)
	assert.Equal(ts.T(), foundOrder, order)
}

func (ts *LomsTestSuite) TestGetOrderByIdNotFound() {
	ctx := context.Background()
	var orderId int64 = 1

	ts.mockOrdersProvider.EXPECT().GetOrderByOrderId(ctx, orderId).Return(models.Order{}, ErrOrderNotFound)

	order, err := ts.loms.OrderInfo(ctx, orderId)
	assert.ErrorIs(ts.T(), err, ErrOrderNotFound)
	assert.Empty(ts.T(), order)
}
