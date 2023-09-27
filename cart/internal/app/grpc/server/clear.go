package server

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	cartpb "route256/cart/api/proto"
)

func (s *GrpcServer) Clear(ctx context.Context, request *cartpb.ClearRequest) (*emptypb.Empty, error) {
	err := s.service.DeleteItemsByUserId(ctx, request.User)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &emptypb.Empty{}, nil
}
