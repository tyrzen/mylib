-- +goose Up
-- +goose StatementBegin
CREATE TABLE sessions
(
    id            UUID PRIMARY KEY                     DEFAULT GEN_RANDOM_UUID(),
    reader_id     UUID REFERENCES readers (id),
    refresh_token TEXT                        NOT NULL,
    expires_at    TIMESTAMP WITHOUT TIME ZONE NOT NULL,
    created_at    TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT NOW()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE sessions;
-- +goose StatementEnd
