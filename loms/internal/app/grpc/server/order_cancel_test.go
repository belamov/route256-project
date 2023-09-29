package server

import (
	"context"

	"route256/loms/internal/app/grpc/pb"
	"route256/loms/internal/app/services"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *LomsGrpcServerTestSuite) TestOrderCancel() {
	request := pb.OrderCancelRequest{
		OrderId: 1,
	}

	s.mockService.EXPECT().OrderCancel(gomock.Any(), request.OrderId).Times(1).Return(nil)

	_, err := s.client.OrderCancel(context.Background(), &request)
	assert.NoError(s.T(), err)
}

func (s *LomsGrpcServerTestSuite) TestOrderCancelOrderNotFound() {
	request := pb.OrderCancelRequest{
		OrderId: 1,
	}

	s.mockService.EXPECT().OrderCancel(gomock.Any(), request.OrderId).Times(1).Return(services.ErrOrderNotFound)

	response, err := s.client.OrderCancel(context.Background(), &request)
	require.Error(s.T(), err)
	assert.Nil(s.T(), response)

	grpcErr, ok := status.FromError(err)
	require.True(s.T(), ok)

	assert.Equal(s.T(), codes.NotFound, grpcErr.Code())
}
