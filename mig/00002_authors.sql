-- +goose Up
-- +goose StatementBegin
CREATE TABLE authors
(
    id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    first_name TEXT NOT NULL,
    last_name  TEXT NOT NULL,
    created_at TIMESTAMP WITHOUT TIME ZONE DEFAULT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE authors
-- +goose StatementEnd
