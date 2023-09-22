package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"route256/cart/internal/app/models"
	"route256/cart/internal/app/services"
)

type ItemAddRequest struct {
	UserId int64
	Sku    uint32
	Count  uint64
}

func (h *Handler) AddItem(w http.ResponseWriter, r *http.Request) {
	var req ItemAddRequest
	dec := json.NewDecoder(r.Body)
	if err := dec.Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	item := models.CartItem{
		User:  req.UserId,
		Sku:   req.Sku,
		Count: req.Count,
	}
	err := h.cart.AddItem(r.Context(), item)
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
