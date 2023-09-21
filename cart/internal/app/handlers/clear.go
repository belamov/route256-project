package handlers

import (
	"encoding/json"
	"net/http"
)

type ClearRequest struct {
	UserId int64
}

func (h *Handler) Clear(w http.ResponseWriter, r *http.Request) {
	var req ClearRequest
	dec := json.NewDecoder(r.Body)
	if err := dec.Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err := h.cart.DeleteItemsByUserId(r.Context(), req.UserId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
