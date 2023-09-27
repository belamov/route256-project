package server

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	cartpb "route256/cart/api/proto"
)

func (s *GrpcServer) List(ctx context.Context, request *cartpb.ListRequest) (*cartpb.ListResponse, error) {
	items, totalPrice, err := s.service.GetItemsByUserId(ctx, request.User)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	response := &cartpb.ListResponse{
		TotalPrice: totalPrice,
		Items:      make([]*cartpb.ListItemResponse, 0, len(items)),
	}
	for _, item := range items {
		response.Items = append(response.Items, &cartpb.ListItemResponse{
			Sku:   item.Sku,
			Count: item.Count,
			Name:  item.Name,
			Price: item.Price,
		})
	}

	return response, nil
}
