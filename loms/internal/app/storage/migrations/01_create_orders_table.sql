-- +goose Up
-- +goose StatementBegin
create table orders (
 id bigserial primary key,
 created_at timestamp not null,
 user_id bigint not null,
 status smallint not null
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table orders;
-- +goose StatementEnd
