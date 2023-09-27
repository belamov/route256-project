package server

import (
	"context"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	cartpb "route256/cart/api/proto"
)

func (s *CartGrpcServerTestSuite) TestDeleteItem() {
	request := cartpb.DeleteItemRequest{
		User: 1,
		Sku:  1,
	}

	s.mockService.EXPECT().DeleteItem(gomock.Any(), gomock.Any()).Times(1).Return(nil)

	_, err := s.client.DeleteItem(context.Background(), &request)
	assert.NoError(s.T(), err)
}
