package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/rs/zerolog/log"
)

type StockInfoRequest struct {
	Sku uint32 `json:"sku,omitempty"`
}

type StockInfoResponse struct {
	Count uint64 `json:"count,omitempty"`
}

func (h *Handler) StockInfo(w http.ResponseWriter, r *http.Request) {
	var req StockInfoRequest
	dec := json.NewDecoder(r.Body)
	if err := dec.Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	count, err := h.loms.StockInfo(r.Context(), req.Sku)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resp := StockInfoResponse{
		Count: count,
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
