package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"route256/loms/internal/app/models"
	"route256/loms/internal/app/services"

	"github.com/rs/zerolog/log"
)

type CreateOrderRequest struct {
	Items  []OrderItemRequest `json:"items,omitempty"`
	UserId int64              `json:"userId,omitempty"`
}

type OrderItemRequest struct {
	Sku   int32  `json:"sku,omitempty"`
	Count uint16 `json:"count,omitempty"`
}

type CreateOrderResponse struct {
	OrderId int64 `json:"orderID,omitempty"`
}

func (h *Handler) CreateOrder(w http.ResponseWriter, r *http.Request) {
	var req CreateOrderRequest
	dec := json.NewDecoder(r.Body)
	if err := dec.Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	orderItems := make([]models.OrderItem, 0, len(req.Items))
	for _, reqItem := range req.Items {
		orderItems = append(orderItems, models.OrderItem{
			Sku:   reqItem.Sku,
			Count: reqItem.Count,
		})
	}

	order, err := h.loms.CreateOrder(r.Context(), req.UserId, orderItems)
	if errors.Is(err, services.ErrInsufficientStocks) {
		http.Error(w, err.Error(), http.StatusPreconditionFailed)
		return
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resp := CreateOrderResponse{
		OrderId: order.Id,
	}

	w.Header().Set("Content-Type", "application/json")
	enc := json.NewEncoder(w)
	if err := enc.Encode(resp); err != nil {
		log.Err(err).Msg("error encoding response")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
