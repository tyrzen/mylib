-- +goose Up
-- +goose StatementBegin
CREATE TABLE books
(
    id        UUID PRIMARY KEY DEFAULT GEN_RANDOM_UUID(),
    author_id UUID        REFERENCES authors (id) ON DELETE SET NULL,
    title     TEXT UNIQUE NOT NULL,
    genre     TEXT        NOT NULL,
    rate      INTEGER     NOT NULL     DEFAULT 0,
    size      INTEGER     NOT NULL,
    year      INTEGER     NOT NULL
);

CREATE INDEX IF NOT EXISTS books_id_idx ON books USING BTREE(id);

CREATE INDEX IF NOT EXISTS books_title_idx ON books USING GIN(to_tsvector('simple', title));
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS books_title_idx;

DROP INDEX IF EXISTS books_id_idx;

DROP TABLE books CASCADE;
-- +goose StatementEnd
