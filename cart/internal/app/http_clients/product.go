package http_clients

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"route256/cart/internal/app/domain/models"
	"route256/cart/internal/app/domain/services"

	"github.com/rs/zerolog/log"
)

type productHttpClient struct {
	client     *http.Client
	serviceUrl string
}

func NewProductHttpClient(serviceUrl string) services.ProductService {
	// TODO: make configurable via env
	defaultTimeout := time.Second * 3
	httpClient := &http.Client{
		Timeout: defaultTimeout,
	}
	return &productHttpClient{
		client:     httpClient,
		serviceUrl: serviceUrl,
	}
}

type GetProductRequest struct {
	Token string `json:"token,omitempty"`
	Sku   uint32 `json:"sku,omitempty"`
}

type GetProductResponse struct {
	Name  string `json:"name"`
	Price uint32 `json:"price"`
}

func (p *productHttpClient) GetProduct(ctx context.Context, sku uint32) (models.CartItemInfo, error) {
	token, ok := ctx.Value("products_token").(string)
	if !ok {
		return models.CartItemInfo{}, errors.New("cant parse products_token from context")
	}

	requestData := GetProductRequest{
		Token: token,
		Sku:   sku,
	}

	buf := new(bytes.Buffer)
	err := json.NewEncoder(buf).Encode(requestData)
	if err != nil {
		return models.CartItemInfo{}, fmt.Errorf("cannot encode response: %w", err)
	}

	if err != nil {
		return models.CartItemInfo{}, fmt.Errorf("cannot create request: %w", err)
	}

	response, err := p.client.Post(
		p.serviceUrl,
		"application/json",
		buf,
	)
	if err != nil {
		log.Error().Err(err).Msg("failed executing http request")
		return models.CartItemInfo{}, err
	}

	if response.StatusCode == http.StatusNotFound {
		return models.CartItemInfo{}, services.ErrSkuInvalid
	}

	respBody, err := io.ReadAll(response.Body)
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(response.Body)
	if err != nil {
		log.Error().Err(err).Msg("failed reading response body")
		return models.CartItemInfo{}, err
	}

	if err != nil {
		log.Error().Err(err).Msg("failed executing http request")
		return models.CartItemInfo{}, err
	}

	var getProductResponse GetProductResponse
	err = json.Unmarshal(respBody, &getProductResponse)
	if err != nil {
		log.Error().Err(err).Bytes("response", respBody).Msg("unexpected response from product service")
		return models.CartItemInfo{}, err
	}

	productInfo := models.CartItemInfo{
		Name:  getProductResponse.Name,
		Sku:   sku,
		Price: getProductResponse.Price,
	}

	return productInfo, nil
}
