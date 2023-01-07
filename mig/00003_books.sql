-- +goose Up
-- +goose StatementBegin
CREATE TABLE books
(
    id         UUID PRIMARY KEY            DEFAULT gen_random_uuid(),
    author_id  UUID REFERENCES authors (id) ON DELETE SET NULL,
    title      TEXT NOT NULL,
    genre      TEXT,
    rate       INTEGER,
    created_at TIMESTAMP WITHOUT TIME ZONE DEFAULT NOW()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE books;
-- +goose StatementEnd
