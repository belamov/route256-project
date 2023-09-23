package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"

	"route256/loms/internal/app/domain/models"
	"route256/loms/internal/app/domain/services"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func (s *HandlersTestSuite) TestHandler_CreateOrder() {
	var userId int64 = 1
	req := CreateOrderRequest{
		UserId: userId,
		Items: []OrderItemRequest{
			{
				Sku:   1,
				Count: 1,
			}, {
				Sku:   2,
				Count: 2,
			},
		},
	}

	body, err := json.Marshal(req)
	require.NoError(s.T(), err)

	createdOrder := models.Order{
		Id: 1,
		Items: []models.OrderItem{
			{
				Sku:   1,
				Count: 1,
			},
			{
				Sku:   2,
				Count: 2,
			},
		},
		Status: models.OrderStatusAwaitingPayment,
	}
	s.mockService.EXPECT().OrderCreate(gomock.Any(), userId, gomock.Any()).Times(1).Return(createdOrder, nil)

	result, response := s.testRequest(
		http.MethodPost,
		"/order/create",
		string(body),
		nil,
	)
	_ = result.Body.Close()

	assert.Equal(s.T(), http.StatusOK, result.StatusCode)

	var resp CreateOrderResponse
	dec := json.NewDecoder(bytes.NewReader([]byte(response)))
	err = dec.Decode(&resp)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), createdOrder.Id, resp.OrderId)
}

func (s *HandlersTestSuite) TestHandler_CreateOrderInsufficientStocks() {
	var userId int64 = 1
	req := CreateOrderRequest{
		UserId: userId,
		Items: []OrderItemRequest{
			{
				Sku:   1,
				Count: 1,
			}, {
				Sku:   2,
				Count: 2,
			},
		},
	}

	body, err := json.Marshal(req)
	require.NoError(s.T(), err)

	createdOrder := models.Order{}
	s.mockService.EXPECT().OrderCreate(gomock.Any(), userId, gomock.Any()).Times(1).Return(createdOrder, services.ErrInsufficientStocks)

	result, _ := s.testRequest(
		http.MethodPost,
		"/order/create",
		string(body),
		nil,
	)
	_ = result.Body.Close()

	assert.Equal(s.T(), http.StatusPreconditionFailed, result.StatusCode)
}
