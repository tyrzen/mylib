-- +goose Up
-- +goose StatementBegin
CREATE TABLE readers
(
    id       UUID NOT NULL PRIMARY KEY,
    first_name TEXT DEFAULT NULL,
    last_name TEXT DEFAULT NULL,
    email    TEXT NOT NULL UNIQUE,
    password TEXT NOT NULL,
    born_at  TIMESTAMP WITHOUT TIME ZONE DEFAULT NULL
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE readers;
-- +goose StatementEnd
