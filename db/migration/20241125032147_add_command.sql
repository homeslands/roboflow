-- +goose Up
-- +goose StatementBegin
CREATE TABLE raybot_commands (
    "id" UUID PRIMARY KEY,
    "raybot_id" UUID NOT NULL,
    "type" TEXT NOT NULL,
    "status" TEXT NOT NULL,
    "input" JSON,
    "output" JSON,
    "created_at" TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    "completed_at" TIMESTAMPTZ
)
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS "raybot_commands";
-- +goose StatementEnd
