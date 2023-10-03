package server

import (
	"context"

	"route256/cart/internal/app/grpc/pb"
	"route256/cart/internal/app/models"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func (s *CartGrpcServerTestSuite) TestList() {
	var userId int64 = 1
	var totalPrice uint32 = 5
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

	request := &pb.ListRequest{User: userId}
	s.mockService.EXPECT().GetItemsByUserId(gomock.Any(), gomock.Any()).Times(1).Return(items, totalPrice, nil)

	response, err := s.client.List(context.Background(), request)
	assert.NoError(s.T(), err)

	assert.Equal(s.T(), totalPrice, response.TotalPrice)

	assert.Equal(s.T(), items[0].Sku, response.Items[0].Sku)
	assert.Equal(s.T(), items[0].Count, response.Items[0].Count)
	assert.Equal(s.T(), items[0].Name, response.Items[0].Name)
	assert.Equal(s.T(), items[0].Price, response.Items[0].Price)

	assert.Equal(s.T(), items[1].Sku, response.Items[1].Sku)
	assert.Equal(s.T(), items[1].Count, response.Items[1].Count)
	assert.Equal(s.T(), items[1].Name, response.Items[1].Name)
	assert.Equal(s.T(), items[1].Price, response.Items[1].Price)
}
