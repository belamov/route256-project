-- name: CreateOrder :one
insert into orders (created_at, user_id, status) values ($1, $2, $3) returning id;

-- name: CreateOrderItems :exec
insert into order_items (order_id, name, sku, count, price) values ($1, $2, $3, $4, $5);

-- name: UpdateOrderStatus :exec
update orders set status = $1 where id = $2;

-- name: GetOrderById :many
select sqlc.embed(o), sqlc.embed(oi) from orders o
    left join order_items oi on o.id = oi.order_id
where o.id = $1;

-- name: GetExpiredOrdersWithStatus :many
select id from orders
where created_at < $1 and status = $2;

-- name: GetBySku :one
select count from stocks where sku = $1;

-- name: ChangeReserveOfSkuByAmount :execresult
update stocks set count = count + $1 where sku = $2 returning id;
