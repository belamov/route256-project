// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.22.0
// source: queries.sql

package queries

import (
	"context"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
)

const changeReserveOfSkuByAmount = `-- name: ChangeReserveOfSkuByAmount :execresult
update stocks set count = count + $1 where sku = $2 returning id
`

type ChangeReserveOfSkuByAmountParams struct {
	Count int64 `json:"count"`
	Sku   int64 `json:"sku"`
}

func (q *Queries) ChangeReserveOfSkuByAmount(ctx context.Context, arg ChangeReserveOfSkuByAmountParams) (pgconn.CommandTag, error) {
	return q.db.Exec(ctx, changeReserveOfSkuByAmount, arg.Count, arg.Sku)
}

const createOrder = `-- name: CreateOrder :one
insert into orders (created_at, user_id, status) values ($1, $2, $3) returning id
`

type CreateOrderParams struct {
	CreatedAt pgtype.Timestamp `json:"created_at"`
	UserID    int64            `json:"user_id"`
	Status    int16            `json:"status"`
}

func (q *Queries) CreateOrder(ctx context.Context, arg CreateOrderParams) (int64, error) {
	row := q.db.QueryRow(ctx, createOrder, arg.CreatedAt, arg.UserID, arg.Status)
	var id int64
	err := row.Scan(&id)
	return id, err
}

const createOrderItems = `-- name: CreateOrderItems :exec
insert into order_items (order_id, name, sku, count, price) values ($1, $2, $3, $4, $5)
`

type CreateOrderItemsParams struct {
	OrderID int64  `json:"order_id"`
	Name    string `json:"name"`
	Sku     int64  `json:"sku"`
	Count   int64  `json:"count"`
	Price   int64  `json:"price"`
}

func (q *Queries) CreateOrderItems(ctx context.Context, arg CreateOrderItemsParams) error {
	_, err := q.db.Exec(ctx, createOrderItems,
		arg.OrderID,
		arg.Name,
		arg.Sku,
		arg.Count,
		arg.Price,
	)
	return err
}

const getBySku = `-- name: GetBySku :one
select count from stocks where sku = $1
`

func (q *Queries) GetBySku(ctx context.Context, sku int64) (int64, error) {
	row := q.db.QueryRow(ctx, getBySku, sku)
	var count int64
	err := row.Scan(&count)
	return count, err
}

const getExpiredOrdersWithStatus = `-- name: GetExpiredOrdersWithStatus :many
select id from orders
where created_at < $1 and status = $2
`

type GetExpiredOrdersWithStatusParams struct {
	CreatedAt pgtype.Timestamp `json:"created_at"`
	Status    int16            `json:"status"`
}

func (q *Queries) GetExpiredOrdersWithStatus(ctx context.Context, arg GetExpiredOrdersWithStatusParams) ([]int64, error) {
	rows, err := q.db.Query(ctx, getExpiredOrdersWithStatus, arg.CreatedAt, arg.Status)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []int64
	for rows.Next() {
		var id int64
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		items = append(items, id)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getOrderById = `-- name: GetOrderById :many
select o.id, o.created_at, o.user_id, o.status, oi.id, oi.order_id, oi.name, oi.sku, oi.count, oi.price from orders o
    left join order_items oi on o.id = oi.order_id
where o.id = $1
`

type GetOrderByIdRow struct {
	Order     Order     `json:"order"`
	OrderItem OrderItem `json:"order_item"`
}

func (q *Queries) GetOrderById(ctx context.Context, id int64) ([]GetOrderByIdRow, error) {
	rows, err := q.db.Query(ctx, getOrderById, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetOrderByIdRow
	for rows.Next() {
		var i GetOrderByIdRow
		if err := rows.Scan(
			&i.Order.ID,
			&i.Order.CreatedAt,
			&i.Order.UserID,
			&i.Order.Status,
			&i.OrderItem.ID,
			&i.OrderItem.OrderID,
			&i.OrderItem.Name,
			&i.OrderItem.Sku,
			&i.OrderItem.Count,
			&i.OrderItem.Price,
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

const updateOrderStatus = `-- name: UpdateOrderStatus :exec
update orders set status = $1 where id = $2
`

type UpdateOrderStatusParams struct {
	Status int16 `json:"status"`
	ID     int64 `json:"id"`
}

func (q *Queries) UpdateOrderStatus(ctx context.Context, arg UpdateOrderStatusParams) error {
	_, err := q.db.Exec(ctx, updateOrderStatus, arg.Status, arg.ID)
	return err
}
