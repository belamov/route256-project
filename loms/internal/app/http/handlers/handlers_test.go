package handlers

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"route256/loms/internal/app/domain/services"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
)

type HandlersTestSuite struct {
	suite.Suite
	mockCtrl    *gomock.Controller
	mockService *services.MockLoms
	ts          *httptest.Server
	r           http.Handler
}

func (s *HandlersTestSuite) SetupSuite() {
	s.mockCtrl = gomock.NewController(Reporter{s.T()})
	s.mockService = services.NewMockLoms(s.mockCtrl)
	s.r = NewRouter(s.mockService)
	s.ts = httptest.NewServer(s.r)
}

func (s *HandlersTestSuite) SetupTest() {
}

func (s *HandlersTestSuite) TearDownTest() {
}

func (s *HandlersTestSuite) TearDownSuite() {
	s.ts.Close()
}

func TestHandlersTestSuite(t *testing.T) {
	suite.Run(t, new(HandlersTestSuite))
}

func (s *HandlersTestSuite) testRequest(method, path string, body string, cookies map[string]string) (*http.Response, string) {
	var err error
	var req *http.Request
	var resp *http.Response
	var respBody []byte

	req, err = http.NewRequest(method, s.ts.URL+path, strings.NewReader(body))
	require.NoError(s.T(), err)

	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	if len(cookies) > 0 {
		for name, value := range cookies {
			req.AddCookie(&http.Cookie{
				Name:     name,
				Value:    value,
				Secure:   true,
				HttpOnly: true,
			})
		}
	}
	resp, err = client.Do(req)
	require.NoError(s.T(), err)

	respBody, err = io.ReadAll(resp.Body)
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)

	require.NoError(s.T(), err)

	return resp, string(bytes.TrimSpace(respBody))
}

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
