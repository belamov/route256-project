-- noinspection SqlInsertValuesForFile
-- name: CreateOrder :exec
insert into orders (created_at, user_id, status) values ($1, $2, $3);