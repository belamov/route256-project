package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"route256/cart/internal/app/models"
	"route256/cart/internal/app/services"
)

func (s *HandlersTestSuite) TestHandler_AddItem() {
	tests := []struct {
		serviceReturn  error
		name           string
		expectedStatus int
	}{
		{name: "no error", serviceReturn: nil, expectedStatus: http.StatusOK},
		{name: "InsufficientStocks", serviceReturn: services.ErrInsufficientStocks, expectedStatus: http.StatusPreconditionFailed},
		{name: "InvalidItem", serviceReturn: services.ErrItemInvalid, expectedStatus: http.StatusBadRequest},
		{name: "InvalidSku", serviceReturn: services.ErrSkuInvalid, expectedStatus: http.StatusBadRequest},
	}
	for _, tt := range tests {
		item := models.CartItem{
			User:  1,
			Sku:   20,
			Count: 5,
		}

		req := ItemAddRequest{
			UserId: item.User,
			Sku:    item.Sku,
			Count:  item.Count,
		}

		body, err := json.Marshal(req)
		require.NoError(s.T(), err)

		s.mockService.EXPECT().AddItem(gomock.Any(), item).Times(1).Return(tt.serviceReturn)

		result, _ := s.testRequest(
			http.MethodPost,
			"/cart/item/add",
			string(body),
			nil,
		)
		_ = result.Body.Close()

		assert.Equal(s.T(), tt.expectedStatus, result.StatusCode)
	}
}
