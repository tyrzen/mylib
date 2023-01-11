-- +goose Up
-- +goose StatementBegin

CREATE TABLE readers
(
    id         UUID PRIMARY KEY            DEFAULT gen_random_uuid(),
    first_name VARCHAR(255)                        DEFAULT NULL,
    last_name  VARCHAR(255)                        DEFAULT NULL,
    email      VARCHAR(255) UNIQUE NOT NULL,
    password   CHAR(60)       NOT NULL,
    created_at TIMESTAMP WITHOUT TIME ZONE DEFAULT NOW()
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE readers;
-- +goose StatementEnd
