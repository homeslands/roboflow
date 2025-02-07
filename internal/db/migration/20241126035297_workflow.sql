-- +goose Up
-- +goose StatementBegin
CREATE TABLE "workflows" (
    "id" UUID NOT NULL PRIMARY KEY,
    "name" TEXT NOT NULL,
    "description" TEXT,
	"is_draft" BOOLEAN NOT NULL DEFAULT TRUE,
	"is_valid" BOOLEAN NOT NULL DEFAULT FALSE,
    "data" JSON NOT NULL DEFAULT '{}',
    "created_at" TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    "updated_at" TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS "workflows";
-- +goose StatementEnd
