package server

import (
	"context"
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	lomspb "route256/loms/api/proto"
	"route256/loms/internal/app/models"
	"route256/loms/internal/app/services"
)

func (s *GrpcServer) OrderCreate(ctx context.Context, request *lomspb.OrderCreateRequest) (*lomspb.OrderCreateResponse, error) {
	orderItems := make([]models.OrderItem, 0, len(request.Items))
	for _, item := range request.Items {
		orderItems = append(orderItems, models.OrderItem{
			Sku:   item.Sku,
			Count: item.Count,
		})
	}

	order, err := s.service.OrderCreate(ctx, request.User, orderItems)
	if errors.Is(err, services.ErrInsufficientStocks) {
		return nil, status.Error(codes.FailedPrecondition, err.Error())
	}
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &lomspb.OrderCreateResponse{OrderId: order.Id}, nil
}
