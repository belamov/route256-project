package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"

	"route256/cart/internal/app/services"

	"github.com/stretchr/testify/assert"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func (s *HandlersTestSuite) TestHandler_Checkout() {
	var userId int64 = 1
	req := ListRequest{
		UserId: userId,
	}

	body, err := json.Marshal(req)
	require.NoError(s.T(), err)

	var orderId int64 = 1000

	s.mockService.EXPECT().Checkout(gomock.Any(), userId).Times(1).Return(orderId, nil)

	result, response := s.testRequest(
		http.MethodPost,
		"/cart/checkout",
		string(body),
		nil,
	)
	_ = result.Body.Close()

	assert.Equal(s.T(), http.StatusOK, result.StatusCode)

	var resp CheckoutResponse
	dec := json.NewDecoder(bytes.NewReader([]byte(response)))
	err = dec.Decode(&resp)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), orderId, resp.OrderId)
}

func (s *HandlersTestSuite) TestHandler_CheckoutEmptyCart() {
	var userId int64 = 1
	req := ListRequest{
		UserId: userId,
	}

	body, err := json.Marshal(req)
	require.NoError(s.T(), err)

	s.mockService.EXPECT().Checkout(gomock.Any(), userId).Times(1).Return(int64(0), services.ErrCartIsEmpty)

	result, _ := s.testRequest(
		http.MethodPost,
		"/cart/checkout",
		string(body),
		nil,
	)
	_ = result.Body.Close()

	assert.Equal(s.T(), http.StatusBadRequest, result.StatusCode)
}
