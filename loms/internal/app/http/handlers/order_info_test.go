package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"

	"route256/loms/internal/app/models"
	"route256/loms/internal/app/services"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func (s *HandlersTestSuite) TestHandler_OrderInfo() {
	var orderId int64 = 1
	var userId int64 = 1
	req := OrderInfoRequest{
		OrderId: orderId,
	}

	body, err := json.Marshal(req)
	require.NoError(s.T(), err)

	foundOrder := models.Order{
		Id: orderId,
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
		UserId: userId,
	}

	s.mockService.EXPECT().OrderInfo(gomock.Any(), orderId).Return(foundOrder, nil)

	result, response := s.testRequest(
		http.MethodPost,
		"/order/info",
		string(body),
		nil,
	)
	_ = result.Body.Close()

	assert.Equal(s.T(), http.StatusOK, result.StatusCode)

	var resp OrderInfoResponse
	dec := json.NewDecoder(bytes.NewReader([]byte(response)))
	err = dec.Decode(&resp)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), foundOrder.Status.String(), resp.Status)
	assert.Equal(s.T(), foundOrder.UserId, resp.User)
}

func (s *HandlersTestSuite) TestHandler_OrderInfoNotFound() {
	var orderId int64 = 1
	req := OrderInfoRequest{
		OrderId: orderId,
	}

	body, err := json.Marshal(req)
	require.NoError(s.T(), err)

	s.mockService.EXPECT().OrderInfo(gomock.Any(), orderId).Return(models.Order{}, services.ErrOrderNotFound)

	result, _ := s.testRequest(
		http.MethodPost,
		"/order/info",
		string(body),
		nil,
	)
	_ = result.Body.Close()

	assert.Equal(s.T(), http.StatusNotFound, result.StatusCode)
}
