-- +goose Up
-- +goose StatementBegin
create table reserves (
 id bigserial primary key,
 created_at timestamp not null,
 sku bigint not null,
 reserved int not null
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table reserves;
-- +goose StatementEnd
