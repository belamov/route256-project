package server

import (
	"context"

	"route256/loms/internal/app/grpc/pb"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func (s *LomsGrpcServerTestSuite) TestStockInfo() {
	request := pb.StockInfoRequest{
		Sku: 1,
	}

	var count uint64 = 10
	s.mockService.EXPECT().StockInfo(gomock.Any(), request.Sku).Times(1).Return(count, nil)

	response, err := s.client.StockInfo(context.Background(), &request)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), count, response.Count)
}
