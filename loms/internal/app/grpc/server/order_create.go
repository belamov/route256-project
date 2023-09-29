package server

import (
	"context"
	"errors"

	"route256/loms/internal/app/grpc/pb"
	"route256/loms/internal/app/models"
	"route256/loms/internal/app/services"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *GrpcServer) OrderCreate(ctx context.Context, request *pb.OrderCreateRequest) (*pb.OrderCreateResponse, error) {
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

	return &pb.OrderCreateResponse{OrderId: order.Id}, nil
}
