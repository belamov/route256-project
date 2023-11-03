package server

import (
	"context"
	"fmt"
	"math/rand"
	"net"
	"net/http"
	"sync"
	"testing"
	"time"

	"route256/loms/internal/pkg/metrics"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/test/bufconn"

	"route256/loms/internal/app/grpc/pb"
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
	client      pb.LomsClient
	conn        *grpc.ClientConn
	mockCtrl    *gomock.Controller
	grpcServer  *grpc.Server
	mockService *services.MockLoms
}

func (s *LomsGrpcServerTestSuite) SetupTest() {
	grpcServer := grpc.NewServer()

	ctrl := gomock.NewController(Reporter{s.T()})

	mockService := services.NewMockLoms(ctrl)

	appServer := NewGRPCServer("", "", mockService, metrics.InitMetrics())

	pb.RegisterLomsServer(grpcServer, appServer)

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

	s.client = pb.NewLomsClient(conn)
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

func (s *LomsGrpcServerTestSuite) TestRunServer() {
	ctx, cancel := context.WithCancel(context.Background())
	wg := &sync.WaitGroup{}
	address := fmt.Sprintf("0.0.0.0:%d", chooseRandomUnusedPort())
	gatewayPort := chooseRandomUnusedPort()
	gatewayAddress := fmt.Sprintf("0.0.0.0:%d", gatewayPort)
	server := NewGRPCServer(address, gatewayAddress, s.mockService, metrics.InitMetrics())

	wg.Add(2)
	go server.Run(ctx, wg)
	go server.RunGateway(ctx, wg)

	conn, err := grpc.DialContext(
		ctx,
		address,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	require.NoError(s.T(), err)
	grpcClient := pb.NewLomsClient(conn)

	waitForGrpcServerStart(grpcClient)

	s.mockService.EXPECT().StockInfo(gomock.Any(), gomock.Any()).Return(uint64(2), nil)
	count, err := grpcClient.StockInfo(context.Background(), &pb.StockInfoRequest{Sku: 20})
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), uint64(2), count.Count)

	waitForHTTPServerStart(gatewayPort)
	resp, err := http.Get(fmt.Sprintf("http://localhost:%d/", gatewayPort))
	_ = resp.Body.Close()
	require.NoError(s.T(), err)
	assert.Equal(s.T(), http.StatusNotFound, resp.StatusCode)

	cancel()
	wg.Wait()
}

func chooseRandomUnusedPort() (port int) {
	for i := 0; i < 10; i++ {
		port = 40000 + int(rand.Int31n(10000))
		if ln, err := net.Listen("tcp", fmt.Sprintf(":%d", port)); err == nil {
			_ = ln.Close()
			break
		}
	}
	return port
}

func waitForGrpcServerStart(client pb.LomsClient) {
	for i := 0; i < 1000; i++ {
		_, err := client.StockInfo(context.Background(), &pb.StockInfoRequest{Sku: 20})
		if err == nil {
			return
		}
		grpcErr, _ := status.FromError(err)
		if grpcErr.Code() != codes.Unavailable {
			return
		}
		time.Sleep(time.Millisecond * 100)
	}
}

func waitForHTTPServerStart(port int) {
	// wait for up to 10 seconds for server to start before returning it
	client := http.Client{Timeout: time.Second}
	defer client.CloseIdleConnections()
	for i := 0; i < 100; i++ {
		if resp, err := client.Get(fmt.Sprintf("http://localhost:%d/ping", port)); err == nil {
			_ = resp.Body.Close()
			return
		}
		time.Sleep(time.Millisecond * 100)
	}
}
