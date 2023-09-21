package handlers

import (
	"encoding/json"
	"net/http"
)

type OrderCancelRequest struct {
	OrderId int64 `json:"orderID,omitempty"`
}

func (h *Handler) OrderCancel(w http.ResponseWriter, r *http.Request) {
	var req OrderCancelRequest
	dec := json.NewDecoder(r.Body)
	if err := dec.Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err := h.loms.OrderCancel(r.Context(), req.OrderId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
