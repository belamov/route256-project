package server

import (
	"context"
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	lomspb "route256/loms/api/proto"
	"route256/loms/internal/app/services"
)

func (s *GrpcServer) OrderPay(ctx context.Context, request *lomspb.OrderPayRequest) (*emptypb.Empty, error) {
	err := s.service.OrderPay(ctx, request.OrderId)
	if errors.Is(err, services.ErrOrderNotFound) {
		return nil, status.Error(codes.NotFound, err.Error())
	}
	if errors.Is(err, services.ErrOrderCancelled) {
		return nil, status.Error(codes.FailedPrecondition, err.Error())
	}
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &emptypb.Empty{}, nil
}
