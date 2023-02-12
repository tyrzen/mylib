-- +goose Up
-- +goose StatementBegin

CREATE TABLE readers
(
    id         UUID PRIMARY KEY            DEFAULT GEN_RANDOM_UUID(),
    first_name VARCHAR(255)                DEFAULT NULL,
    last_name  VARCHAR(255)                DEFAULT NULL,
    email      VARCHAR(255) UNIQUE NOT NULL,
    password   CHAR(60)            NOT NULL,
    role       CHAR(16)                    DEFAULT NULL,
    created_at TIMESTAMP WITHOUT TIME ZONE DEFAULT NOW(),
    deleted_at TIMESTAMP WITHOUT TIME ZONE DEFAULT NULL
);

CREATE INDEX IF NOT EXISTS readers_id_email_idx ON readers using BTREE(id, email) WHERE deleted_at IS NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS readers_id_email_idx;

DROP TABLE readers CASCADE;
-- +goose StatementEnd
