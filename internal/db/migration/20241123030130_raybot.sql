-- +goose Up
-- +goose StatementBegin
CREATE TABLE "raybots" (
    "id" UUID NOT NULL PRIMARY KEY,
    "name" TEXT NOT NULL UNIQUE,
	"control_mode" TEXT NOT NULL,
	"is_online" BOOLEAN NOT NULL,
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
