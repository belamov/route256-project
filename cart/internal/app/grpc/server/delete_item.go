package server

import (
	"context"
	"errors"

	"route256/cart/internal/app/grpc/pb"
	"route256/cart/internal/app/models"
	"route256/cart/internal/app/services"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *GrpcServer) DeleteItem(ctx context.Context, request *pb.DeleteItemRequest) (*emptypb.Empty, error) {
	cartItem := models.CartItem{
		User:  request.User,
		Sku:   request.Sku,
		Count: 1,
	}

	err := s.service.DeleteItem(ctx, cartItem)
	if errors.Is(err, services.ErrItemInvalid) {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &emptypb.Empty{}, nil
}
