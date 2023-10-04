// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.22.0
// source: queries.sql

package queries

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const createOrder = `-- name: CreateOrder :exec
insert into orders (created_at, user_id, status) values ($1, $2, $3)
`

type CreateOrderParams struct {
	CreatedAt pgtype.Timestamp `json:"created_at"`
	UserID    int64            `json:"user_id"`
	Status    int16            `json:"status"`
}

// noinspection SqlInsertValuesForFile
func (q *Queries) CreateOrder(ctx context.Context, arg CreateOrderParams) error {
	_, err := q.db.Exec(ctx, createOrder, arg.CreatedAt, arg.UserID, arg.Status)
	return err
}
