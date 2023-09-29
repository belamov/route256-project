package server

import (
	"context"
	"errors"

	"route256/loms/internal/app/grpc/pb"
	"route256/loms/internal/app/services"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *GrpcServer) OrderInfo(ctx context.Context, request *pb.OrderInfoRequest) (*pb.OrderInfoResponse, error) {
	err := request.Validate()
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	order, err := s.service.OrderInfo(ctx, request.OrderId)
	if errors.Is(err, services.ErrOrderNotFound) {
		return nil, status.Error(codes.NotFound, err.Error())
	}
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	response := pb.OrderInfoResponse{
		Status: order.Status.String(),
		User:   order.UserId,
		Items:  make([]*pb.OrderItemInfoResponse, 0, len(order.Items)),
	}
	for _, item := range order.Items {
		response.Items = append(response.Items, &pb.OrderItemInfoResponse{
			Sku:   item.Sku,
			Count: item.Count,
			Name:  item.Name,
			Price: item.Price,
		})
	}

	return &response, nil
}
