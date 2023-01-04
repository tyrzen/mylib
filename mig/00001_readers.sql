-- +goose Up
-- +goose StatementBegin
CREATE TABLE readers
(
    id         UUID PRIMARY KEY,
    first_name TEXT DEFAULT NULL,
    last_name  TEXT DEFAULT NULL,
    email      TEXT UNIQUE NOT NULL,
    password   TEXT        NOT NULL,
    created_at TIMESTAMP WITHOUT TIME ZONE DEFAULT NOW()
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE readers;
-- +goose StatementEnd
