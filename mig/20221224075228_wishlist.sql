-- +goose Up
-- +goose StatementBegin
CREATE table wishlist_books
(
    reader_id UUID NOT NULL REFERENCES readers (id),
    book_id   UUID NOT NULL REFERENCES books (id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE wishlist_books;
-- +goose StatementEnd
