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

func (s *LomsGrpcServerTestSuite) TestOrderCreate() {
	request := pb.OrderCreateRequest{
		User: 1,
		Items: []*pb.OrderItemCreateRequest{
			{
				Sku:   1,
				Count: 1,
			},
			{
				Sku:   2,
				Count: 2,
			},
		},
	}

	s.mockService.EXPECT().OrderCreate(gomock.Any(), gomock.Any(), gomock.Any()).Times(1).Return(models.Order{
		CreatedAt: time.Now(),
		Items:     []models.OrderItem{},
		Id:        1,
		UserId:    1,
		Status:    models.OrderStatusAwaitingPayment,
	}, nil)

	_, err := s.client.OrderCreate(context.Background(), &request)
	assert.NoError(s.T(), err)
}

func (s *LomsGrpcServerTestSuite) TestOrderCreateInsufficientStocks() {
	request := pb.OrderCreateRequest{
		User: 1,
		Items: []*pb.OrderItemCreateRequest{
			{
				Sku:   1,
				Count: 1,
			},
			{
				Sku:   2,
				Count: 2,
			},
		},
	}

	s.mockService.EXPECT().OrderCreate(gomock.Any(), gomock.Any(), gomock.Any()).Times(1).Return(models.Order{}, services.ErrInsufficientStocks)

	response, err := s.client.OrderCreate(context.Background(), &request)
	require.Error(s.T(), err)
	assert.Nil(s.T(), response)

	grpcErr, ok := status.FromError(err)
	require.True(s.T(), ok)

	assert.Equal(s.T(), codes.FailedPrecondition, grpcErr.Code())
}
