package server

import (
	"context"

	"route256/cart/internal/app/grpc/pb"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *GrpcServer) Clear(ctx context.Context, request *pb.ClearRequest) (*emptypb.Empty, error) {
	err := request.Validate()
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	err = s.service.DeleteItemsByUserId(ctx, request.User)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &emptypb.Empty{}, nil
}
