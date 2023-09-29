package server

import (
	"context"
	"errors"

	"route256/loms/internal/app/grpc/pb"
	"route256/loms/internal/app/services"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *GrpcServer) OrderCancel(ctx context.Context, request *pb.OrderCancelRequest) (*emptypb.Empty, error) {
	err := request.Validate()
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	err = s.service.OrderCancel(ctx, request.OrderId)
	if errors.Is(err, services.ErrOrderNotFound) {
		return nil, status.Error(codes.NotFound, err.Error())
	}
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &emptypb.Empty{}, nil
}
