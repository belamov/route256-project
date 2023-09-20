package handlers

import (
	"encoding/json"
	"net/http"

	"route256/cart/internal/app/models"
)

type ItemDeleteRequest struct {
	UserId int64
	Sku    uint32
}

func (h *Handler) DeleteItem(w http.ResponseWriter, r *http.Request) {
	var req ItemDeleteRequest
	dec := json.NewDecoder(r.Body)
	if err := dec.Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	item := models.CartItem{
		User: req.UserId,
		Sku:  req.Sku,
	}
	err := h.cart.DeleteItem(r.Context(), item)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
