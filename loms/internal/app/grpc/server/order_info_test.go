package server

import (
	"context"
	"time"

	"route256/loms/internal/app/grpc/pb"
	"route256/loms/internal/app/models"
	"route256/loms/internal/app/services"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *LomsGrpcServerTestSuite) TestOrderInfo() {
	request := pb.OrderInfoRequest{
		OrderId: 1,
	}

	order := models.Order{
		CreatedAt: time.Now(),
		Items:     []models.OrderItem{},
		Id:        1,
		UserId:    1,
		Status:    models.OrderStatusAwaitingPayment,
	}
	s.mockService.EXPECT().OrderInfo(gomock.Any(), request.OrderId).Times(1).Return(order, nil)

	response, err := s.client.OrderInfo(context.Background(), &request)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), order.Status.String(), response.Status)
	assert.Equal(s.T(), order.UserId, response.User)
	assert.Equal(s.T(), len(order.Items), len(response.Items))
}

func (s *LomsGrpcServerTestSuite) TestOrderInfoNotFound() {
	request := pb.OrderInfoRequest{
		OrderId: 1,
	}

	order := models.Order{}
	s.mockService.EXPECT().OrderInfo(gomock.Any(), request.OrderId).Times(1).Return(order, services.ErrOrderNotFound)

	response, err := s.client.OrderInfo(context.Background(), &request)
	require.Error(s.T(), err)
	assert.Nil(s.T(), response)

	grpcErr, ok := status.FromError(err)
	require.True(s.T(), ok)

	assert.Equal(s.T(), codes.NotFound, grpcErr.Code())
}
