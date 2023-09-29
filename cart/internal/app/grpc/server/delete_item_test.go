package server

import (
	"context"

	"route256/cart/internal/app/grpc/pb"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func (s *CartGrpcServerTestSuite) TestDeleteItem() {
	request := pb.DeleteItemRequest{
		User: 1,
		Sku:  1,
	}

	s.mockService.EXPECT().DeleteItem(gomock.Any(), gomock.Any()).Times(1).Return(nil)

	_, err := s.client.DeleteItem(context.Background(), &request)
	assert.NoError(s.T(), err)
}
