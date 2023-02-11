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
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE books CASCADE;
-- +goose StatementEnd
