package server

import (
	"context"

	"route256/cart/internal/app/grpc/pb"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func (s *CartGrpcServerTestSuite) TestClear() {
	request := pb.ClearRequest{
		User: 1,
	}

	s.mockService.EXPECT().DeleteItemsByUserId(gomock.Any(), request.User).Times(1).Return(nil)

	_, err := s.client.Clear(context.Background(), &request)
	assert.NoError(s.T(), err)
}
