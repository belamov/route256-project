package server

import (
	"context"
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	lomspb "route256/loms/api/proto"
	"route256/loms/internal/app/services"
)

func (s *GrpcServer) OrderInfo(ctx context.Context, request *lomspb.OrderInfoRequest) (*lomspb.OrderInfoResponse, error) {
	order, err := s.service.OrderInfo(ctx, request.OrderId)
	if errors.Is(err, services.ErrOrderNotFound) {
		return nil, status.Error(codes.NotFound, err.Error())
	}
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	response := lomspb.OrderInfoResponse{
		Status: order.Status.String(),
		User:   order.UserId,
		Items:  make([]*lomspb.OrderItemInfoResponse, 0, len(order.Items)),
	}
	for _, item := range order.Items {
		response.Items = append(response.Items, &lomspb.OrderItemInfoResponse{
			Sku:   item.Sku,
			Count: item.Count,
			Name:  item.Name,
			Price: item.Price,
		})
	}

	return &response, nil
}
