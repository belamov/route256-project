package server

import (
	"context"

	"route256/cart/internal/app/grpc/pb"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *GrpcServer) Checkout(ctx context.Context, request *pb.CheckoutRequest) (*pb.CheckoutResponse, error) {
	orderId, err := s.service.Checkout(ctx, request.User)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	response := pb.CheckoutResponse{OrderID: orderId}

	return &response, nil
}
