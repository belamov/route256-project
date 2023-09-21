package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func (s *HandlersTestSuite) TestHandler_OrderPay() {
	var orderId int64 = 1
	req := OrderPayRequest{
		OrderId: orderId,
	}

	body, err := json.Marshal(req)
	require.NoError(s.T(), err)

	s.mockService.EXPECT().OrderPay(gomock.Any(), orderId).Return(nil)

	result, _ := s.testRequest(
		http.MethodPost,
		"/order/pay",
		string(body),
		nil,
	)
	_ = result.Body.Close()

	assert.Equal(s.T(), http.StatusOK, result.StatusCode)
}
