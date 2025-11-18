-- +goose Up
-- +goose StatementBegin
create table metrics (
    id text primary key,
    mtype text not null,
    delta bigint,
    value double precision,
    created_at timestamp default current_timestamp not null,
    updated_at timestamp default current_timestamp not null
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table metrics;
-- +goose StatementEnd
