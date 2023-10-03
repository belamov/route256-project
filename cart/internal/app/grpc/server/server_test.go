package server

import (
	"context"
	"fmt"
	"net"
	"testing"

	"route256/cart/internal/app/grpc/pb"
	"route256/cart/internal/app/services"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"
)

type Reporter struct {
	T *testing.T
}

// ensure Reporter implements gomock.TestReporter.
var _ gomock.TestReporter = Reporter{}

// Errorf is equivalent testing.T.Errorf.
func (r Reporter) Errorf(format string, args ...interface{}) {
	r.T.Errorf(format, args...)
}

// Fatalf crashes the program with a panic to allow users to diagnose
// missing expects.
func (r Reporter) Fatalf(format string, args ...interface{}) {
	panic(fmt.Sprintf(format, args...))
}

type CartGrpcServerTestSuite struct {
	suite.Suite
	client      pb.CartClient
	conn        *grpc.ClientConn
	mockCtrl    *gomock.Controller
	grpcServer  *grpc.Server
	mockService *services.MockCart
}

func (s *CartGrpcServerTestSuite) SetupTest() {
	grpcServer := grpc.NewServer()

	ctrl := gomock.NewController(Reporter{s.T()})

	mockService := services.NewMockCart(ctrl)

	appServer := NewGRPCServer("", "", mockService)

	pb.RegisterCartServer(grpcServer, appServer)

	s.mockCtrl = ctrl
	s.grpcServer = grpcServer

	listener := bufconn.Listen(1024 * 1024)

	dialer := func(context.Context, string) (net.Conn, error) {
		return listener.Dial()
	}

	conn, err := grpc.DialContext(
		context.Background(),
		"",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithContextDialer(dialer),
	)
	require.NoError(s.T(), err)

	go func(s *CartGrpcServerTestSuite) {
		errServe := s.grpcServer.Serve(listener)
		require.NoError(s.T(), errServe)
	}(s)

	s.client = pb.NewCartClient(conn)
	s.conn = conn
	s.mockService = mockService
}

func (s *CartGrpcServerTestSuite) TearDownTest() {
	s.mockCtrl.Finish()
	s.grpcServer.Stop()
}

func (s *CartGrpcServerTestSuite) TearDownSuite() {
	err := s.conn.Close()
	require.NoError(s.T(), err)
}

func TestCartGrpcServerTestSuite(t *testing.T) {
	suite.Run(t, new(CartGrpcServerTestSuite))
}
