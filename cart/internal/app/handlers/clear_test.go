package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/stretchr/testify/assert"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func (s *HandlersTestSuite) TestHandler_Clear() {
	var userId int64 = 1
	req := ClearRequest{
		UserId: userId,
	}

	body, err := json.Marshal(req)
	require.NoError(s.T(), err)

	s.mockService.EXPECT().DeleteItemsByUserId(gomock.Any(), userId).Times(1).Return(nil)

	result, _ := s.testRequest(
		http.MethodPost,
		"/cart/clear",
		string(body),
		nil,
	)
	_ = result.Body.Close()

	assert.Equal(s.T(), http.StatusOK, result.StatusCode)
}
