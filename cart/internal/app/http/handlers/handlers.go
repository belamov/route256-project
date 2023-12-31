package handlers

import (
	"net/http"

	"route256/cart/internal/app/services"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type Handler struct {
	cart services.Cart
}

func NewHandler(service services.Cart) *Handler {
	return &Handler{
		cart: service,
	}
}

func NewRouter(cart services.Cart) http.Handler {
	h := NewHandler(cart)

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Heartbeat("/ping"))

	r.Post("/cart/item/add", h.AddItem)
	r.Post("/cart/item/delete", h.DeleteItem)
	r.Post("/cart/list", h.List)
	r.Post("/cart/clear", h.Clear)
	r.Post("/cart/checkout", h.Checkout)

	return r
}
