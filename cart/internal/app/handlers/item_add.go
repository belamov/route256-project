package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"route256/cart/internal/app/domain/models"
	"route256/cart/internal/app/domain/services"
)

type ItemAddRequest struct {
	UserId int64  `json:"user,omitempty"`
	Sku    uint32 `json:"sku,omitempty"`
	Count  uint64 `json:"count,omitempty"`
}

func (r *ItemAddRequest) Validate() error {
	if r.UserId == 0 {
		return errors.New("user is required")
	}

	if r.Sku == 0 {
		return errors.New("sku is required")
	}

	if r.Count <= 0 {
		return errors.New("count must be greater than 0")
	}

	return nil
}

func (h *Handler) AddItem(w http.ResponseWriter, r *http.Request) {
	var req ItemAddRequest
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

	item := models.CartItem{
		User:  req.UserId,
		Sku:   req.Sku,
		Count: req.Count,
	}
	err = h.cart.AddItem(r.Context(), item)
	if errors.Is(err, services.ErrItemInvalid) || errors.Is(err, services.ErrSkuInvalid) {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if errors.Is(err, services.ErrInsufficientStocks) {
		http.Error(w, err.Error(), http.StatusPreconditionFailed)
		return
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
