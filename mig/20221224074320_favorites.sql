-- +goose Up
-- +goose StatementBegin
CREATE table favorites
(
    reader_id UUID NOT NULL REFERENCES readers (id),
    book_id   UUID NOT NULL REFERENCES books (id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE favorites;
-- +goose StatementEnd
