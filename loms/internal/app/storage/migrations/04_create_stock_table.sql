-- +goose Up
-- +goose StatementBegin
create table stocks (
 id bigserial primary key,
 sku bigint not null,
 count bigint not null
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table stocks;
-- +goose StatementEnd
