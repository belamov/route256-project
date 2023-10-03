package server

import (
	"context"

	"route256/cart/internal/app/grpc/pb"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func (s *CartGrpcServerTestSuite) TestCheckout() {
	request := pb.CheckoutRequest{
		User: 1,
	}

	var orderId int64 = 2

	s.mockService.EXPECT().Checkout(gomock.Any(), request.User).Times(1).Return(orderId, nil)

	response, err := s.client.Checkout(context.Background(), &request)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), orderId, response.OrderID)
}
