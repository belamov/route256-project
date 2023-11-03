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

	"route256/cart/internal/pkg/metrics"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

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

	appServer := NewGRPCServer("", "", mockService, metrics.InitMetrics())

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

func (s *CartGrpcServerTestSuite) TestRunServer() {
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
	grpcClient := pb.NewCartClient(conn)

	waitForGrpcServerStart(grpcClient)

	_, err = grpcClient.AddItem(context.Background(), &pb.AddItemRequest{
		User: 1,
		Item: nil,
	})
	grpcErr, _ := status.FromError(err)
	assert.Equal(s.T(), codes.InvalidArgument, grpcErr.Code())

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

func waitForGrpcServerStart(client pb.CartClient) {
	for i := 0; i < 1000; i++ {
		_, err := client.List(context.Background(), &pb.ListRequest{User: 1})
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
