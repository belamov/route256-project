package server

import (
	"context"

	"route256/cart/internal/app/grpc/pb"
	"route256/cart/internal/app/services"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *CartGrpcServerTestSuite) TestAddItem() {
	request := pb.AddItemRequest{
		User: 1,
		Item: &pb.CartItemAddRequest{
			User:  1,
			Sku:   1,
			Count: 1,
		},
	}

	s.mockService.EXPECT().AddItem(gomock.Any(), gomock.Any()).Times(1).Return(nil)

	_, err := s.client.AddItem(context.Background(), &request)
	assert.NoError(s.T(), err)
}

func (s *CartGrpcServerTestSuite) TestAddItemInsufficientStocks() {
	request := pb.AddItemRequest{
		User: 1,
		Item: &pb.CartItemAddRequest{
			User:  1,
			Sku:   1,
			Count: 1,
		},
	}

	s.mockService.EXPECT().AddItem(gomock.Any(), gomock.Any()).Times(1).Return(services.ErrInsufficientStocks)

	response, err := s.client.AddItem(context.Background(), &request)
	require.Error(s.T(), err)
	assert.Nil(s.T(), response)

	grpcErr, ok := status.FromError(err)
	require.True(s.T(), ok)

	assert.Equal(s.T(), codes.FailedPrecondition, grpcErr.Code())
}
