-- +goose Up
-- +goose StatementBegin
CREATE TABLE raybots (
    "id" UUID PRIMARY KEY,
    "name" TEXT NOT NULL UNIQUE,
    "token" TEXT NOT NULL,
    "status" TEXT NOT NULL,
    "created_at" TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    "updated_at" TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TRIGGER update_raybots_updated_at_trigger
    BEFORE UPDATE ON raybots
    FOR EACH ROW
    EXECUTE PROCEDURE update_row_updated_at_column();
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS "raybots";
-- +goose StatementEnd
