package server

import (
	"context"
	"fmt"
	"net"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"

	lomspb "route256/loms/api/proto"
	"route256/loms/internal/app/services"
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

type LomsGrpcServerTestSuite struct {
	suite.Suite
	client      lomspb.LomsClient
	conn        *grpc.ClientConn
	mockCtrl    *gomock.Controller
	grpcServer  *grpc.Server
	mockService *services.MockLoms
}

func (s *LomsGrpcServerTestSuite) SetupTest() {
	grpcServer := grpc.NewServer()

	ctrl := gomock.NewController(Reporter{s.T()})

	mockService := services.NewMockLoms(ctrl)

	appServer := NewGRPCServer("", mockService)

	lomspb.RegisterLomsServer(grpcServer, appServer)

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

	go func(s *LomsGrpcServerTestSuite) {
		errServe := s.grpcServer.Serve(listener)
		require.NoError(s.T(), errServe)
	}(s)

	s.client = lomspb.NewLomsClient(conn)
	s.conn = conn
	s.mockService = mockService
}

func (s *LomsGrpcServerTestSuite) TearDownTest() {
	s.mockCtrl.Finish()
	s.grpcServer.Stop()
}

func (s *LomsGrpcServerTestSuite) TearDownSuite() {
	err := s.conn.Close()
	require.NoError(s.T(), err)
}

func TestLomsGrpcServerTestSuite(t *testing.T) {
	suite.Run(t, new(LomsGrpcServerTestSuite))
}
