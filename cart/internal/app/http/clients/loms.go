package clients

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"route256/cart/internal/app/domain/models"
	"route256/cart/internal/app/domain/services"

	"github.com/rs/zerolog/log"
)

type lomsHttpClient struct {
	client     *http.Client
	serviceUrl string
}

func NewLomsHttpClient(serviceUrl string) services.LomsService {
	httpClient := &http.Client{
		Timeout: time.Second * 3,
	}
	return &lomsHttpClient{
		client:     httpClient,
		serviceUrl: serviceUrl,
	}
}

type StockInfoRequest struct {
	Sku uint32 `json:"sku,omitempty"`
}

type StockInfoResponse struct {
	Count uint64 `json:"count,omitempty"`
}

func (l *lomsHttpClient) GetStocksInfo(ctx context.Context, sku uint32) (uint64, error) {
	requestData := StockInfoRequest{
		Sku: sku,
	}

	buf := new(bytes.Buffer)
	err := json.NewEncoder(buf).Encode(requestData)
	if err != nil {
		return 0, fmt.Errorf("cannot encode response: %w", err)
	}

	if err != nil {
		return 0, fmt.Errorf("cannot create request: %w", err)
	}

	endpoint := &url.URL{
		Scheme: "https",
		Host:   l.serviceUrl,
		Path:   "/stock/info",
	}
	response, err := l.client.Post(
		endpoint.String(),
		"application/json",
		buf,
	)
	if err != nil {
		log.Error().Err(err).Msg("failed executing http request")
		return 0, err
	}

	if response.StatusCode == http.StatusNotFound {
		return 0, services.ErrSkuInvalid
	}

	respBody, err := io.ReadAll(response.Body)
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(response.Body)
	if err != nil {
		log.Error().Err(err).Msg("failed reading response body")
		return 0, err
	}

	if err != nil {
		log.Error().Err(err).Msg("failed executing http request")
		return 0, err
	}

	var stockInfoResponse StockInfoResponse
	err = json.Unmarshal(respBody, &stockInfoResponse)
	if err != nil {
		log.Error().Err(err).Bytes("response", respBody).Msg("unexpected response from loms service")
		return 0, err
	}

	return stockInfoResponse.Count, nil
}

type CreateOrderRequest struct {
	Items  []OrderItemRequest `json:"items,omitempty"`
	UserId int64              `json:"userId,omitempty"`
}

type OrderItemRequest struct {
	Sku   uint32 `json:"sku,omitempty"`
	Count uint64 `json:"count,omitempty"`
}

type CreateOrderResponse struct {
	OrderId int64 `json:"orderID,omitempty"`
}

func (l *lomsHttpClient) CreateOrder(ctx context.Context, userId int64, items []models.CartItem) (int64, error) {
	requestData := CreateOrderRequest{
		Items:  make([]OrderItemRequest, 0, len(items)),
		UserId: 0,
	}
	for _, cartItem := range items {
		requestData.Items = append(requestData.Items, OrderItemRequest{
			Sku:   cartItem.Sku,
			Count: cartItem.Count,
		})
	}

	buf := new(bytes.Buffer)
	err := json.NewEncoder(buf).Encode(requestData)
	if err != nil {
		return 0, fmt.Errorf("cannot encode response: %w", err)
	}

	if err != nil {
		return 0, fmt.Errorf("cannot create request: %w", err)
	}

	endpoint := &url.URL{
		Scheme: "https",
		Host:   l.serviceUrl,
		Path:   "/order/create",
	}
	response, err := l.client.Post(
		endpoint.String(),
		"application/json",
		buf,
	)
	if err != nil {
		log.Error().Err(err).Msg("failed executing http request")
		return 0, err
	}

	respBody, err := io.ReadAll(response.Body)
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(response.Body)
	if err != nil {
		log.Error().Err(err).Msg("failed reading response body")
		return 0, err
	}

	if err != nil {
		log.Error().Err(err).Msg("failed executing http request")
		return 0, err
	}

	var createOrderResponse CreateOrderResponse
	err = json.Unmarshal(respBody, &createOrderResponse)
	if err != nil {
		log.Error().Err(err).Bytes("response", respBody).Msg("unexpected response from loms service")
		return 0, err
	}

	return createOrderResponse.OrderId, nil
}
