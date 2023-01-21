-- +goose Up
-- +goose StatementBegin
CREATE TABLE books
(
    id         UUID PRIMARY KEY            DEFAULT GEN_RANDOM_UUID(),
    author_id  UUID REFERENCES authors (id) ON DELETE SET NULL,
    title      TEXT NOT NULL,
    genre      TEXT,
    rate       INTEGER,
    created_at TIMESTAMP WITHOUT TIME ZONE DEFAULT NOW()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE books CASCADE;
-- +goose StatementEnd
