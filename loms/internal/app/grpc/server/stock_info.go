package server

import (
	"context"

	"route256/loms/internal/app/grpc/pb"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *GrpcServer) StockInfo(ctx context.Context, request *pb.StockInfoRequest) (*pb.StockInfoResponse, error) {
	stockAvailable, err := s.service.StockInfo(ctx, request.Sku)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.StockInfoResponse{
		Count: stockAvailable,
	}, nil
}
