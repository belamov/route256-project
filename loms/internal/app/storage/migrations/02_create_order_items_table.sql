-- +goose Up
-- +goose StatementBegin
create table order_items (
 id bigserial primary key,
 order_id bigint not null references orders(id),
 name varchar not null,
 sku bigint not null,
 count bigint not null,
 price bigint not null
);


-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table order_items;
-- +goose StatementEnd
