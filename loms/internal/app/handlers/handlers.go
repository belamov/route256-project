package handlers

import (
	"net/http"

	"route256/loms/internal/app/domain/services"

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

	r.Post("/order/create", h.OrderCreate)
	r.Post("/order/info", h.OrderInfo)
	r.Post("/order/pay", h.OrderPay)
	r.Post("/order/cancel", h.OrderCancel)
	r.Post("/stock/info", h.StockInfo)

	return r
}
