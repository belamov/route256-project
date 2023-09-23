package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"

	"route256/cart/internal/app/domain/models"

	"github.com/stretchr/testify/assert"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func (s *HandlersTestSuite) TestHandler_List() {
	var userId int64 = 1
	req := ListRequest{
		UserId: userId,
	}

	body, err := json.Marshal(req)
	require.NoError(s.T(), err)

	var totalPrice uint32 = 10
	items := []models.CartItemWithInfo{
		{
			User:  userId,
			Sku:   1,
			Count: 1,
			Name:  "1",
			Price: 1,
		},
		{
			User:  userId,
			Sku:   2,
			Count: 2,
			Name:  "2",
			Price: 2,
		},
	}
	s.mockService.EXPECT().GetItemsByUserId(gomock.Any(), userId).Times(1).Return(items, totalPrice, nil)

	result, response := s.testRequest(
		http.MethodPost,
		"/cart/list",
		string(body),
		nil,
	)
	_ = result.Body.Close()

	var resp ListResponse
	dec := json.NewDecoder(bytes.NewReader([]byte(response)))
	err = dec.Decode(&resp)
	assert.NoError(s.T(), err)
}
