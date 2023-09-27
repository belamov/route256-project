package server

import (
	"context"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	cartpb "route256/cart/api/proto"
)

func (s *CartGrpcServerTestSuite) TestClear() {
	request := cartpb.ClearRequest{
		User: 1,
	}

	s.mockService.EXPECT().DeleteItemsByUserId(gomock.Any(), request.User).Times(1).Return(nil)

	_, err := s.client.Clear(context.Background(), &request)
	assert.NoError(s.T(), err)
}
