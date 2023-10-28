-- +goose Up
-- +goose StatementBegin
create table outbox (
 id bigserial primary key,
 topic varchar not null,
 data bytea not null,
 key varchar not null,
 sent_at timestamp,
 error_message varchar,
 retry_count smallint,
 locked_by varchar,
 locked_at timestamp,
 created_at timestamp not null
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table outbox;
-- +goose StatementEnd
