-- +goose Up
-- +goose StatementBegin
CREATE TABLE authors
(
    id         UUID PRIMARY KEY DEFAULT GEN_RANDOM_UUID(),
    first_name VARCHAR(255) NOT NULL,
    last_name  VARCHAR(255) NOT NULL,
    created_at TIMESTAMP WITHOUT TIME ZONE DEFAULT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE authors CASCADE;
-- +goose StatementEnd
