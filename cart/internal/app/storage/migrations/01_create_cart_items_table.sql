-- +goose Up
-- +goose StatementBegin
create table cart_items (
 id bigserial primary key,
 sku bigint not null,
 count bigint not null,
 user_id bigint not null
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table cart_items;
-- +goose StatementEnd
