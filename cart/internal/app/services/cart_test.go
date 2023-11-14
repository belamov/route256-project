package services

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"testing"
	"time"

	"go.opentelemetry.io/otel/trace"

	"route256/cart/internal/app/models"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
)

type CartTestSuite struct {
	suite.Suite
	mockCtrl           *gomock.Controller
	mockCartProvider   *MockCartProvider
	mockLomsService    *MockLomsService
	mockProductService *MockProductService
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
	ts.mockCartProvider = NewMockCartProvider(ts.mockCtrl)
	ts.mockLomsService = NewMockLomsService(ts.mockCtrl)
	ts.mockProductService = NewMockProductService(ts.mockCtrl)
	tracer := trace.NewNoopTracerProvider().Tracer("mock")
	ts.cart = NewCartService(ts.mockProductService, ts.mockLomsService, ts.mockCartProvider, NoopCache{}, tracer)
}

func TestCartTestSuite(t *testing.T) {
	suite.Run(t, new(CartTestSuite))
}

func (ts *CartTestSuite) TestAddItem() {
	ctx := context.Background()
	item := models.CartItem{
		User:  1,
		Sku:   1,
		Count: 1,
	}
	ts.mockLomsService.EXPECT().GetStocksInfo(ctx, item.Sku).Times(1).Return(item.Count+1, nil)
	ts.mockProductService.EXPECT().GetProduct(ctx, item.Sku).Times(1).Return(models.CartItemInfo{}, nil)
	ts.mockCartProvider.EXPECT().SaveItem(ctx, item).Times(1).Return(nil)
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
	ts.mockLomsService.EXPECT().GetStocksInfo(ctx, item.Sku).Times(1).Return(item.Count+1, nil)
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
	ts.mockLomsService.EXPECT().GetStocksInfo(ctx, item.Sku).Times(1).Return(item.Count+1, nil)
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
	ts.mockLomsService.EXPECT().GetStocksInfo(ctx, item.Sku).Times(1).Return(item.Count-1, nil)
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

	ts.mockCartProvider.EXPECT().DeleteItem(ctx, item).Return(nil)
	err := ts.cart.DeleteItem(ctx, item)
	assert.NoError(ts.T(), err)
}

func (ts *CartTestSuite) TestGetItemsByUserIdNoUser() {
	ctx := context.Background()

	items, total, err := ts.cart.GetItemsByUserId(ctx, 0)
	assert.Empty(ts.T(), items)
	assert.Empty(ts.T(), total)
	assert.Error(ts.T(), err)
}

func (ts *CartTestSuite) TestGetItemsByUserId() {
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
	ts.mockCartProvider.EXPECT().GetItemsByUserId(gomock.Any(), userId).Return(cartItems, nil)
	for i, item := range cartItems {
		ts.mockProductService.EXPECT().
			GetProduct(gomock.Any(), item.Sku).
			Return(models.CartItemInfo{
				Sku:   item.Sku,
				Name:  strconv.Itoa(i),
				Price: uint32(item.Count),
			}, nil)
	}
	items, total, err := ts.cart.GetItemsByUserId(ctx, userId)
	assert.NoError(ts.T(), err)
	assert.Len(ts.T(), items, len(cartItems))
	assert.Equal(ts.T(), uint32(1+2*2+3*3), total)
	for i, item := range cartItems {
		assert.Equal(ts.T(), item.Sku, items[i].Sku)
		assert.Equal(ts.T(), item.Count, items[i].Count)
		assert.Equal(ts.T(), strconv.Itoa(i), items[i].Name)
		assert.Equal(ts.T(), uint32(item.Count), items[i].Price)
	}
}

func (ts *CartTestSuite) TestGetItemsByUserIdWithError() {
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
	returnedErr := errors.New("some error")

	ts.mockCartProvider.EXPECT().GetItemsByUserId(gomock.Any(), userId).Return(cartItems, nil)
	ts.mockProductService.EXPECT().GetProduct(gomock.Any(), cartItems[0].Sku).Return(models.CartItemInfo{
		Sku:   cartItems[0].Sku,
		Name:  "name",
		Price: 100,
	}, nil).AnyTimes()
	ts.mockProductService.EXPECT().GetProduct(gomock.Any(), cartItems[1].Sku).Return(models.CartItemInfo{}, returnedErr).Times(1)
	ts.mockProductService.EXPECT().GetProduct(gomock.Any(), cartItems[2].Sku).Return(models.CartItemInfo{
		Sku:   cartItems[0].Sku,
		Name:  "name",
		Price: 100,
	}, nil).AnyTimes()
	items, total, err := ts.cart.GetItemsByUserId(ctx, userId)
	assert.ErrorIs(ts.T(), err, returnedErr)
	assert.Empty(ts.T(), items)
	assert.Empty(ts.T(), total)
}

func (ts *CartTestSuite) TestCheckout() {
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
	ts.mockCartProvider.EXPECT().GetItemsByUserId(gomock.Any(), userId).Return(cartItems, nil)
	ts.mockCartProvider.EXPECT().DeleteItemsByUserId(gomock.Any(), userId).Return(nil)
	orderId := int64(1000)
	ts.mockLomsService.EXPECT().
		CreateOrder(gomock.Any(), userId, gomock.Any()).
		Return(orderId, nil)

	returnedOrderId, err := ts.cart.Checkout(ctx, userId)
	assert.NoError(ts.T(), err)
	assert.Equal(ts.T(), orderId, returnedOrderId)
}

func (ts *CartTestSuite) TestCheckoutEmptyCart() {
	ctx := context.Background()
	var userId int64 = 1
	var cartItems []models.CartItem
	ts.mockCartProvider.EXPECT().GetItemsByUserId(gomock.Any(), userId).Return(cartItems, nil)

	returnedOrderId, err := ts.cart.Checkout(ctx, userId)
	assert.ErrorIs(ts.T(), err, ErrCartIsEmpty)
	assert.Empty(ts.T(), returnedOrderId)
}

func (ts *CartTestSuite) TestCheckoutNoUser() {
	ctx := context.Background()
	var userId int64 = 0
	orderId, err := ts.cart.Checkout(ctx, userId)
	assert.Error(ts.T(), err)
	assert.Equal(ts.T(), int64(0), orderId)
}

func (ts *CartTestSuite) TestDeleteItemsByUserId() {
	ctx := context.Background()
	var userId int64 = 1
	ts.mockCartProvider.EXPECT().DeleteItemsByUserId(ctx, userId).Return(nil)
	err := ts.cart.DeleteItemsByUserId(ctx, userId)
	assert.NoError(ts.T(), err)
}

func (ts *CartTestSuite) TestDeleteItemsByUserIdNoUser() {
	ctx := context.Background()
	var userId int64 = 0
	err := ts.cart.DeleteItemsByUserId(ctx, userId)
	assert.Error(ts.T(), err)
}

type mockProductService struct{}

func (m mockProductService) GetProduct(ctx context.Context, sku uint32) (models.CartItemInfo, error) {
	time.Sleep(time.Millisecond * 10)
	return models.CartItemInfo{
		Name:  "some name",
		Sku:   sku,
		Price: 100,
	}, nil
}

// go test -bench=. -benchmem -cpuprofile=../../../testdata/list_bench/base/cpu.out -memprofile=../../../testdata/list_bench/base/mem.out -o ../../../testdata/list_bench/base/ -count=10 . > ../../../testdata/list_bench/base/result.txt
// go test -bench=. -benchmem -cpuprofile=../../../testdata/list_bench/concurrent/cpu.out -memprofile=../../../testdata/list_bench/concurrent/mem.out -o ../../../testdata/list_bench/concurrent/ -count=10 . > ../../../testdata/list_bench/concurrent/result.txt
// go install golang.org/x/perf/cmd/benchstat@latest
// benchstat base/result.txt concurrent/result.txt
func BenchmarkCart_GetItemsByUserId(b *testing.B) {
	ctrl := gomock.NewController(b)
	mockCartStorage := NewMockCartProvider(ctrl)
	mockLomsService := NewMockLomsService(ctrl)

	var userId int64 = 50

	// чтобы не делать настоящих запросов, замокаем продукт сервис с обычным слипом в 10мс, как будто бы мы получаем ответ за 10мс
	service := NewCartService(mockProductService{}, mockLomsService, mockCartStorage, NoopCache{}, nil)

	cartWith1item := make([]models.CartItem, 1)
	for j := 0; j < len(cartWith1item); j++ {
		cartWith1item[j] = models.CartItem{
			User:  userId,
			Sku:   123,
			Count: uint64(j),
		}
	}

	b.Run(fmt.Sprintf("GetCartItems Cart List With %d items", len(cartWith1item)), func(b *testing.B) {
		mockCartStorage.EXPECT().GetItemsByUserId(gomock.Any(), userId).Return(cartWith1item, nil).Times(b.N)
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_, _, err := service.GetItemsByUserId(context.Background(), userId)
			if err != nil {
				b.Fatal(err)
			}
		}
	})

	cartWith2items := make([]models.CartItem, 2)
	for j := 0; j < len(cartWith2items); j++ {
		cartWith2items[j] = models.CartItem{
			User:  userId,
			Sku:   123,
			Count: uint64(j),
		}
	}

	b.Run(fmt.Sprintf("GetCartItems Cart List With %d items", len(cartWith2items)), func(b *testing.B) {
		mockCartStorage.EXPECT().GetItemsByUserId(gomock.Any(), userId).Return(cartWith2items, nil).Times(b.N)
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_, _, err := service.GetItemsByUserId(context.Background(), userId)
			if err != nil {
				b.Fatal(err)
			}
		}
	})

	cartWith5items := make([]models.CartItem, 5)
	for j := 0; j < len(cartWith5items); j++ {
		cartWith5items[j] = models.CartItem{
			User:  userId,
			Sku:   123,
			Count: uint64(j),
		}
	}

	b.Run(fmt.Sprintf("GetCartItems Cart List With %d items", len(cartWith5items)), func(b *testing.B) {
		mockCartStorage.EXPECT().GetItemsByUserId(gomock.Any(), userId).Return(cartWith5items, nil).Times(b.N)
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_, _, err := service.GetItemsByUserId(context.Background(), userId)
			if err != nil {
				b.Fatal(err)
			}
		}
	})

	cartWith10items := make([]models.CartItem, 10)
	for j := 0; j < len(cartWith10items); j++ {
		cartWith10items[j] = models.CartItem{
			User:  userId,
			Sku:   123,
			Count: uint64(j),
		}
	}

	b.Run(fmt.Sprintf("GetCartItems Cart List With %d items", len(cartWith10items)), func(b *testing.B) {
		mockCartStorage.EXPECT().GetItemsByUserId(gomock.Any(), userId).Return(cartWith10items, nil).Times(b.N)
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_, _, err := service.GetItemsByUserId(context.Background(), userId)
			if err != nil {
				b.Fatal(err)
			}
		}
	})

	cartWith20items := make([]models.CartItem, 20)
	for j := 0; j < len(cartWith20items); j++ {
		cartWith20items[j] = models.CartItem{
			User:  userId,
			Sku:   123,
			Count: uint64(j),
		}
	}

	b.Run(fmt.Sprintf("GetCartItems Cart List With %d items", len(cartWith20items)), func(b *testing.B) {
		mockCartStorage.EXPECT().GetItemsByUserId(gomock.Any(), userId).Return(cartWith20items, nil).Times(b.N)
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_, _, err := service.GetItemsByUserId(context.Background(), userId)
			if err != nil {
				b.Fatal(err)
			}
		}
	})
}
