-- noinspection SqlInsertValuesForFile
-- name: SaveCartItem :exec
insert into cart_items (sku, count, user_id) values ($1, $2, $3);

-- name: DeleteCartItem :exec
delete from cart_items where user_id=$1 and sku=$2;

-- name: DeleteCartItemsByUserId :exec
delete from cart_items where user_id=$1;

-- name: GetCartItemsByUserId :many
select * from cart_items where user_id=$1;
