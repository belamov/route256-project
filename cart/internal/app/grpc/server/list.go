package server

import (
	"context"

	"route256/cart/internal/app/grpc/pb"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *GrpcServer) List(ctx context.Context, request *pb.ListRequest) (*pb.ListResponse, error) {
	err := request.Validate()
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	items, totalPrice, err := s.service.GetItemsByUserId(ctx, request.User)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	response := &pb.ListResponse{
		TotalPrice: totalPrice,
		Items:      make([]*pb.ListItemResponse, 0, len(items)),
	}
	for _, item := range items {
		response.Items = append(response.Items, &pb.ListItemResponse{
			Sku:   item.Sku,
			Count: item.Count,
			Name:  item.Name,
			Price: item.Price,
		})
	}

	return response, nil
}
