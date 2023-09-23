package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
)

type ClearRequest struct {
	UserId int64 `json:"user,omitempty"`
}

func (r *ClearRequest) Validate() error {
	if r.UserId == 0 {
		return errors.New("user is required")
	}

	return nil
}

func (h *Handler) Clear(w http.ResponseWriter, r *http.Request) {
	var req ClearRequest
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

	err = h.cart.DeleteItemsByUserId(r.Context(), req.UserId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
