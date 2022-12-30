-- +goose Up
-- +goose StatementBegin
CREATE table wishlist
(
    reader_id UUID NOT NULL REFERENCES readers (id),
    book_id   UUID NOT NULL REFERENCES books (id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE wishlist;
-- +goose StatementEnd
