package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func (s *HandlersTestSuite) TestHandler_OrderCancel() {
	var orderId int64 = 1
	req := OrderCancelRequest{
		OrderId: orderId,
	}

	body, err := json.Marshal(req)
	require.NoError(s.T(), err)

	s.mockService.EXPECT().OrderCancel(gomock.Any(), orderId).Return(nil)

	result, _ := s.testRequest(
		http.MethodPost,
		"/order/cancel",
		string(body),
		nil,
	)
	_ = result.Body.Close()

	assert.Equal(s.T(), http.StatusOK, result.StatusCode)
}
