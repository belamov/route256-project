package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"route256/cart/internal/app/domain/models"
)

type ItemDeleteRequest struct {
	UserId int64  `json:"user,omitempty"`
	Sku    uint32 `json:"sku,omitempty"`
}

func (r *ItemDeleteRequest) Validate() error {
	if r.UserId == 0 {
		return errors.New("user is required")
	}

	if r.Sku == 0 {
		return errors.New("sku is required")
	}

	return nil
}

func (h *Handler) DeleteItem(w http.ResponseWriter, r *http.Request) {
	var req ItemDeleteRequest
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
		User: req.UserId,
		Sku:  req.Sku,
	}
	err = h.cart.DeleteItem(r.Context(), item)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
