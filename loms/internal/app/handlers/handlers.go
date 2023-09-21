package handlers

import (
	"net/http"

	"route256/loms/internal/app/services"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type Handler struct {
	loms services.Loms
}

func NewHandler(service services.Loms) *Handler {
	return &Handler{
		loms: service,
	}
}

func NewRouter(loms services.Loms) http.Handler {
	h := NewHandler(loms)

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Heartbeat("/ping"))

	r.Post("/order/create", h.CreateOrder)
	r.Post("/order/info", h.OrderInfo)

	return r
}
