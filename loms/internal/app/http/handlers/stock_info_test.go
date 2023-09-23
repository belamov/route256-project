package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/stretchr/testify/assert"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func (s *HandlersTestSuite) TestHandler_StockInfo() {
	var sku uint32 = 1
	req := StockInfoRequest{
		Sku: sku,
	}

	var skuCount uint64 = 1
	s.mockService.EXPECT().StockInfo(gomock.Any(), sku).Times(1).Return(skuCount, nil)

	body, err := json.Marshal(req)
	require.NoError(s.T(), err)

	result, response := s.testRequest(
		http.MethodPost,
		"/stock/info",
		string(body),
		nil,
	)
	_ = result.Body.Close()

	assert.Equal(s.T(), http.StatusOK, result.StatusCode)

	var resp StockInfoResponse
	dec := json.NewDecoder(bytes.NewReader([]byte(response)))
	err = dec.Decode(&resp)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), skuCount, resp.Count)
}
