-- +goose Up
-- +goose StatementBegin
CREATE table favorites
(
    reader_id UUID NOT NULL REFERENCES readers (id) ON DELETE CASCADE,
    book_id   UUID NOT NULL REFERENCES books (id) ON DELETE CASCADE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE favorites;
-- +goose StatementEnd
