-- +goose Up
-- +goose StatementBegin

CREATE TABLE readers
(
    id         UUID PRIMARY KEY            DEFAULT GEN_RANDOM_UUID(),
    first_name VARCHAR(255)                DEFAULT NULL,
    last_name  VARCHAR(255)                DEFAULT NULL,
    email      VARCHAR(255) UNIQUE NOT NULL,
    password   CHAR(60)            NOT NULL,
    is_admin   BOOL                        DEFAULT FALSE,
    created_at TIMESTAMP WITHOUT TIME ZONE DEFAULT NOW()
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE readers CASCADE;
-- +goose StatementEnd
