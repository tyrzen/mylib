
-- +goose Up
-- +goose StatementBegin
CREATE TABLE books
(
    id        UUID PRIMARY KEY,
    author_id UUID,
    title     TEXT NOT NULL,
    genre     TEXT,
    rate      INTEGER,
    wrote_at  TIMESTAMP WITHOUT TIME ZONE DEFAULT NOW(),

    CONSTRAINT fk_author
        FOREIGN KEY (author_id)
            REFERENCES authors(id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE books;
-- +goose StatementEnd
