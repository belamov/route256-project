package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"route256/loms/internal/app/services"

	"github.com/rs/zerolog/log"
)

type OrderInfoRequest struct {
	OrderId int64 `json:"orderID,omitempty"`
}

func (r *OrderInfoRequest) Validate() error {
	if r.OrderId == 0 {
		return errors.New("orderID is required")
	}

	return nil
}

type OrderInfoResponse struct {
	Status string                  `json:"status,omitempty"`
	Items  []OrderItemInfoResponse `json:"items,omitempty"`
	User   int64                   `json:"user,omitempty"`
}

type OrderItemInfoResponse struct {
	Sku   uint32 `json:"sku,omitempty"`
	Count uint64 `json:"count,omitempty"`
}

func (h *Handler) OrderInfo(w http.ResponseWriter, r *http.Request) {
	var req OrderInfoRequest
	dec := json.NewDecoder(r.Body)
	if err := dec.Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err := req.Validate()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	order, err := h.loms.OrderInfo(r.Context(), req.OrderId)
	if errors.Is(err, services.ErrOrderNotFound) {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resp := OrderInfoResponse{
		Status: order.Status.String(),
		User:   order.UserId,
		Items:  make([]OrderItemInfoResponse, 0, len(order.Items)),
	}
	for _, orderItem := range order.Items {
		resp.Items = append(resp.Items, OrderItemInfoResponse{
			Sku:   orderItem.Sku,
			Count: orderItem.Count,
		})
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
