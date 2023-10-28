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

const getLockedUnsentMessage = `-- name: GetLockedUnsentMessage :many
select id, topic, data, key, sent_at, error_message, retry_count, locked_by, locked_at, created_at from outbox where locked_by=$1 and sent_at is null
`

func (q *Queries) GetLockedUnsentMessage(ctx context.Context, lockedBy pgtype.Text) ([]Outbox, error) {
	rows, err := q.db.Query(ctx, getLockedUnsentMessage, lockedBy)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Outbox
	for rows.Next() {
		var i Outbox
		if err := rows.Scan(
			&i.ID,
			&i.Topic,
			&i.Data,
			&i.Key,
			&i.SentAt,
			&i.ErrorMessage,
			&i.RetryCount,
			&i.LockedBy,
			&i.LockedAt,
			&i.CreatedAt,
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

const lockUnsentMessages = `-- name: LockUnsentMessages :exec
update outbox set locked_by=$1, locked_at=$2 where locked_by is null and locked_at is null and sent_at is null
`

type LockUnsentMessagesParams struct {
	LockedBy pgtype.Text      `json:"locked_by"`
	LockedAt pgtype.Timestamp `json:"locked_at"`
}

func (q *Queries) LockUnsentMessages(ctx context.Context, arg LockUnsentMessagesParams) error {
	_, err := q.db.Exec(ctx, lockUnsentMessages, arg.LockedBy, arg.LockedAt)
	return err
}

const saveOutboxMessage = `-- name: SaveOutboxMessage :exec
insert into outbox (key, topic, data) values ($1, $2, $3)
`

type SaveOutboxMessageParams struct {
	Key   string `json:"key"`
	Topic string `json:"topic"`
	Data  []byte `json:"data"`
}

func (q *Queries) SaveOutboxMessage(ctx context.Context, arg SaveOutboxMessageParams) error {
	_, err := q.db.Exec(ctx, saveOutboxMessage, arg.Key, arg.Topic, arg.Data)
	return err
}

const setMessageFailed = `-- name: SetMessageFailed :exec
update outbox set error_message = $1, retry_count = retry_count+1 where id=$2
`

type SetMessageFailedParams struct {
	ErrorMessage pgtype.Text `json:"error_message"`
	ID           int64       `json:"id"`
}

func (q *Queries) SetMessageFailed(ctx context.Context, arg SetMessageFailedParams) error {
	_, err := q.db.Exec(ctx, setMessageFailed, arg.ErrorMessage, arg.ID)
	return err
}

const setMessageSent = `-- name: SetMessageSent :exec
update outbox set sent_at = $1 where id=$2
`

type SetMessageSentParams struct {
	SentAt pgtype.Timestamp `json:"sent_at"`
	ID     int64            `json:"id"`
}

func (q *Queries) SetMessageSent(ctx context.Context, arg SetMessageSentParams) error {
	_, err := q.db.Exec(ctx, setMessageSent, arg.SentAt, arg.ID)
	return err
}

const unlockUnsentMessages = `-- name: UnlockUnsentMessages :exec
update outbox set locked_by=$1, locked_at=$2 where sent_at is null
`

type UnlockUnsentMessagesParams struct {
	LockedBy pgtype.Text      `json:"locked_by"`
	LockedAt pgtype.Timestamp `json:"locked_at"`
}

func (q *Queries) UnlockUnsentMessages(ctx context.Context, arg UnlockUnsentMessagesParams) error {
	_, err := q.db.Exec(ctx, unlockUnsentMessages, arg.LockedBy, arg.LockedAt)
	return err
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
