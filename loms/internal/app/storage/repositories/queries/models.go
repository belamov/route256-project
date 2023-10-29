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

type Outbox struct {
	ID           int64            `json:"id"`
	Destination  string           `json:"destination"`
	Data         []byte           `json:"data"`
	Key          string           `json:"key"`
	SentAt       pgtype.Timestamp `json:"sent_at"`
	ErrorMessage pgtype.Text      `json:"error_message"`
	RetryCount   pgtype.Int2      `json:"retry_count"`
	LockedBy     pgtype.Text      `json:"locked_by"`
	LockedAt     pgtype.Timestamp `json:"locked_at"`
	CreatedAt    pgtype.Timestamp `json:"created_at"`
}
