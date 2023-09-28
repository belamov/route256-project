package server

import (
	"context"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	lomspb "route256/loms/api/proto"
)

func (s *LomsGrpcServerTestSuite) TestStockInfo() {
	request := lomspb.StockInfoRequest{
		Sku: 1,
	}

	var count uint64 = 10
	s.mockService.EXPECT().StockInfo(gomock.Any(), request.Sku).Times(1).Return(count, nil)

	response, err := s.client.StockInfo(context.Background(), &request)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), count, response.Count)
}
