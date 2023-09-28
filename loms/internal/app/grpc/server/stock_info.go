package server

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	lomspb "route256/loms/api/proto"
)

func (s *GrpcServer) StockInfo(ctx context.Context, request *lomspb.StockInfoRequest) (*lomspb.StockInfoResponse, error) {
	stockAvailable, err := s.service.StockInfo(ctx, request.Sku)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &lomspb.StockInfoResponse{
		Count: stockAvailable,
	}, nil
}
