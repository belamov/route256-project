package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/rs/zerolog/log"
)

type ListRequest struct {
	UserId int64
}

func (r *ListRequest) Validate() error {
	if r.UserId == 0 {
		return errors.New("user is required")
	}

	return nil
}

type ListResponse struct {
	Items      []ListItemResponse `json:"items,omitempty"`
	TotalPrice uint32             `json:"total_price,omitempty"`
}

type ListItemResponse struct {
	Name  string `json:"name,omitempty"`
	Sku   uint32 `json:"sku,omitempty"`
	Price uint32 `json:"price,omitempty"`
	Count uint64 `json:"count,omitempty"`
}

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	var req ListRequest
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

	items, totalPrice, err := h.cart.GetItemsByUserId(r.Context(), req.UserId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	res := ListResponse{
		TotalPrice: totalPrice,
		Items:      make([]ListItemResponse, 0, len(items)),
	}
	res.TotalPrice = totalPrice
	for _, item := range items {
		res.Items = append(res.Items, ListItemResponse{
			Sku:   item.Sku,
			Count: item.Count,
			Name:  item.Name,
			Price: item.Price,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	enc := json.NewEncoder(w)
	if err := enc.Encode(res); err != nil {
		log.Err(err).Msg("error encoding response")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
