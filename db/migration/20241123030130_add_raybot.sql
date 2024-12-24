-- +goose Up
-- +goose StatementBegin
CREATE TABLE raybots (
    "id" UUID PRIMARY KEY,
    "name" TEXT NOT NULL UNIQUE,
    "token" TEXT NOT NULL,
    "status" TEXT NOT NULL,
    "ip_address" TEXT,
    "last_connected_at" TIMESTAMPTZ,
    "created_at" TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    "updated_at" TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS "raybots";
-- +goose StatementEnd
