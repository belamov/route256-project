package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
)

type OrderCancelRequest struct {
	OrderId int64 `json:"orderID,omitempty"`
}

func (r *OrderCancelRequest) Validate() error {
	if r.OrderId == 0 {
		return errors.New("orderID is required")
	}

	return nil
}

func (h *Handler) OrderCancel(w http.ResponseWriter, r *http.Request) {
	var req OrderCancelRequest
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

	err = h.loms.OrderCancel(r.Context(), req.OrderId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
