-- +goose Up
-- +goose StatementBegin
CREATE table wishlists
(
    id         UUID PRIMARY KEY            DEFAULT GEN_RANDOM_UUID(),
    reader_id  UUID NOT NULL REFERENCES readers (id) ON DELETE CASCADE,
    book_id    UUID NOT NULL REFERENCES books (id) ON DELETE CASCADE,
    created_at TIMESTAMP WITHOUT TIME ZONE DEFAULT NOW()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE wishlists;
-- +goose StatementEnd
