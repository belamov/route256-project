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

func (s *LomsGrpcServerTestSuite) TestOrderPay() {
	request := pb.OrderPayRequest{
		OrderId: 1,
	}

	s.mockService.EXPECT().OrderPay(gomock.Any(), request.OrderId).Times(1).Return(nil)

	_, err := s.client.OrderPay(context.Background(), &request)
	assert.NoError(s.T(), err)
}

func (s *LomsGrpcServerTestSuite) TestOrderPayCancelledOrder() {
	request := pb.OrderPayRequest{
		OrderId: 1,
	}

	s.mockService.EXPECT().OrderPay(gomock.Any(), request.OrderId).Times(1).Return(services.ErrOrderCancelled)

	response, err := s.client.OrderPay(context.Background(), &request)
	require.Error(s.T(), err)
	assert.Nil(s.T(), response)

	grpcErr, ok := status.FromError(err)
	require.True(s.T(), ok)

	assert.Equal(s.T(), codes.FailedPrecondition, grpcErr.Code())
}

func (s *LomsGrpcServerTestSuite) TestOrderPayOrderNotFound() {
	request := pb.OrderPayRequest{
		OrderId: 1,
	}

	s.mockService.EXPECT().OrderPay(gomock.Any(), request.OrderId).Times(1).Return(services.ErrOrderNotFound)

	response, err := s.client.OrderPay(context.Background(), &request)
	require.Error(s.T(), err)
	assert.Nil(s.T(), response)

	grpcErr, ok := status.FromError(err)
	require.True(s.T(), ok)

	assert.Equal(s.T(), codes.NotFound, grpcErr.Code())
}
