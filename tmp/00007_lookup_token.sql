-- +goose Up
-- +goose StatementBegin
CREATE FUNCTION del_old_tokens() RETURNS trigger
    LANGUAGE plpgsql
AS
$$
BEGIN
    DELETE
    FROM tokens
    WHERE expires_at < NOW();
    RETURN NULL;
END;
$$;

CREATE TRIGGER remove_old_tokens
    AFTER INSERT
    ON tokens
EXECUTE PROCEDURE del_old_tokens();
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TRIGGER remove_old_tokens ON tokens;

DROP FUNCTION del_old_tokens;
-- +goose StatementEnd
