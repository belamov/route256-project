-- +goose Up
-- +goose StatementBegin
create table stocks (
 id bigserial primary key,
 sku bigint not null unique,
 count bigint not null,
 constraint count_nonnegative check (stocks.count >= 0)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table stocks;
-- +goose StatementEnd
