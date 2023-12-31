// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.22.0
// source: queries.sql

package queries

import (
	"context"
)

const deleteCartItem = `-- name: DeleteCartItem :exec
delete from cart_items where user_id=$1 and sku=$2
`

type DeleteCartItemParams struct {
	UserID int64 `json:"user_id"`
	Sku    int64 `json:"sku"`
}

func (q *Queries) DeleteCartItem(ctx context.Context, arg DeleteCartItemParams) error {
	_, err := q.db.Exec(ctx, deleteCartItem, arg.UserID, arg.Sku)
	return err
}

const deleteCartItemsByUserId = `-- name: DeleteCartItemsByUserId :exec
delete from cart_items where user_id=$1
`

func (q *Queries) DeleteCartItemsByUserId(ctx context.Context, userID int64) error {
	_, err := q.db.Exec(ctx, deleteCartItemsByUserId, userID)
	return err
}

const getCartItemsByUserId = `-- name: GetCartItemsByUserId :many
select id, sku, count, user_id from cart_items where user_id=$1
`

func (q *Queries) GetCartItemsByUserId(ctx context.Context, userID int64) ([]CartItem, error) {
	rows, err := q.db.Query(ctx, getCartItemsByUserId, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []CartItem
	for rows.Next() {
		var i CartItem
		if err := rows.Scan(
			&i.ID,
			&i.Sku,
			&i.Count,
			&i.UserID,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const saveCartItem = `-- name: SaveCartItem :exec
insert into cart_items (sku, count, user_id) values ($1, $2, $3)
`

type SaveCartItemParams struct {
	Sku    int64 `json:"sku"`
	Count  int64 `json:"count"`
	UserID int64 `json:"user_id"`
}

// noinspection SqlInsertValuesForFile
func (q *Queries) SaveCartItem(ctx context.Context, arg SaveCartItemParams) error {
	_, err := q.db.Exec(ctx, saveCartItem, arg.Sku, arg.Count, arg.UserID)
	return err
}
