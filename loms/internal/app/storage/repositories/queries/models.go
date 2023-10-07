// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.22.0

package queries

import (
	"github.com/jackc/pgx/v5/pgtype"
)

type Order struct {
	ID        int64            `json:"id"`
	CreatedAt pgtype.Timestamp `json:"created_at"`
	UserID    int64            `json:"user_id"`
	Status    int16            `json:"status"`
}

type OrderItem struct {
	ID      int64  `json:"id"`
	OrderID int64  `json:"order_id"`
	Name    string `json:"name"`
	Sku     int64  `json:"sku"`
	Count   int64  `json:"count"`
	Price   int64  `json:"price"`
}