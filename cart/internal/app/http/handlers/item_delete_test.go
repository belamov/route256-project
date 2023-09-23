package handlers

import (
	"encoding/json"
	"net/http"

	"route256/cart/internal/app/domain/models"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func (s *HandlersTestSuite) TestHandler_DeleteItem() {
	tests := []struct {
		serviceReturn  error
		name           string
		expectedStatus int
	}{
		{name: "no error", serviceReturn: nil, expectedStatus: http.StatusOK},
	}
	for _, tt := range tests {
		item := models.CartItem{
			User: 1,
			Sku:  20,
		}

		req := ItemDeleteRequest{
			UserId: item.User,
			Sku:    item.Sku,
		}

		body, err := json.Marshal(req)
		require.NoError(s.T(), err)

		s.mockService.EXPECT().DeleteItem(gomock.Any(), item).Times(1).Return(tt.serviceReturn)

		result, _ := s.testRequest(
			http.MethodPost,
			"/cart/item/delete",
			string(body),
			nil,
		)
		_ = result.Body.Close()

		assert.Equal(s.T(), tt.expectedStatus, result.StatusCode)
	}
}
