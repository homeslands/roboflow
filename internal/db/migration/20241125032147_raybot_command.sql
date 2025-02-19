-- +goose Up
-- +goose StatementBegin
CREATE TABLE "raybot_commands" (
    "id" UUID NOT NULL PRIMARY KEY,
    "raybot_id" UUID NOT NULL,
    "type" TEXT NOT NULL,
    "status" TEXT NOT NULL,
    "inputs" JSON NOT NULL DEFAULT '{}',
    "outputs" JSON NOT NULL DEFAULT '{}',
    "error" TEXT,
    "created_at" TIMESTAMPTZ NOT NULL DEFAULT NOW(),
	"updated_at" TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    "completed_at" TIMESTAMPTZ,

	FOREIGN KEY("raybot_id") REFERENCES "raybots"("id") ON DELETE CASCADE
)
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS "raybot_commands";
-- +goose StatementEnd
