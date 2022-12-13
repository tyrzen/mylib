-- +goose Up
-- +goose StatementBegin
CREATE TABLE authors (
                         id UUID PRIMARY KEY,
                         first_name TEXT NOT NULL,
                         last_name TEXT NOT NULL,
                         born_at TIMESTAMP WITHOUT TIME ZONE DEFAULT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE authors
-- +goose StatementEnd
