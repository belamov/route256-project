package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"route256/cart/internal/app/services"

	"github.com/rs/zerolog/log"
)

type CheckoutRequest struct {
	UserId int64 `json:"user,omitempty"`
}

func (r *CheckoutRequest) Validate() error {
	if r.UserId == 0 {
		return errors.New("user is required")
	}

	return nil
}

type CheckoutResponse struct {
	OrderId int64 `json:"orderId,omitempty"`
}

func (h *Handler) Checkout(w http.ResponseWriter, r *http.Request) {
	var req CheckoutRequest
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

	orderId, err := h.cart.Checkout(r.Context(), req.UserId)
	if errors.Is(err, services.ErrCartIsEmpty) {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resp := CheckoutResponse{
		OrderId: orderId,
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
