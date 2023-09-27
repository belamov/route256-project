package server

import (
	"context"
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	cartpb "route256/cart/api/proto"
	"route256/cart/internal/app/models"
	"route256/cart/internal/app/services"
)

func (s *GrpcServer) AddItem(ctx context.Context, request *cartpb.AddItemRequest) (*emptypb.Empty, error) {
	cartItem := models.CartItem{
		User:  request.User,
		Sku:   request.Item.Sku,
		Count: request.Item.Count,
	}

	err := s.service.AddItem(ctx, cartItem)
	if errors.Is(err, services.ErrInsufficientStocks) {
		return nil, status.Error(codes.FailedPrecondition, err.Error())
	}
	if errors.Is(err, services.ErrItemInvalid) || errors.Is(err, services.ErrSkuInvalid) {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &emptypb.Empty{}, nil
}
