package server

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	cartpb "route256/cart/api/proto"
)

func (s *GrpcServer) Checkout(ctx context.Context, request *cartpb.CheckoutRequest) (*cartpb.CheckoutResponse, error) {
	orderId, err := s.service.Checkout(ctx, request.User)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	response := cartpb.CheckoutResponse{OrderID: orderId}

	return &response, nil
}
