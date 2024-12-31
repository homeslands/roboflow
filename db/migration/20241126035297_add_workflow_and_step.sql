-- +goose Up
-- +goose StatementBegin
CREATE TABLE workflows (
    "id" UUID PRIMARY KEY,
    "name" TEXT NOT NULL,
    "description" TEXT,
    "definition" JSON NOT NULL DEFAULT '{}',
    "created_at" TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    "updated_at" TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE workflow_executions (
    "id" UUID PRIMARY KEY,
    "workflow_id" UUID NOT NULL,
    "status" TEXT NOT NULL,
    "env" JSON NOT NULL DEFAULT '{}',
    "definition" JSON NOT NULL DEFAULT '{}',
    "created_at" TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    "started_at" TIMESTAMPTZ,
    "completed_at" TIMESTAMPTZ,
    FOREIGN KEY ("workflow_id") REFERENCES "workflows"("id") ON DELETE CASCADE
);

CREATE TABLE steps (
    "id" UUID PRIMARY KEY,
    "workflow_execution_id" UUID NOT NULL,
	"env" JSON NOT NULL DEFAULT '{}',
    "node" JSON NOT NULL DEFAULT '{}',
    "status" TEXT NOT NULL,
    "started_at" TIMESTAMPTZ,
    "completed_at" TIMESTAMPTZ,
    FOREIGN KEY ("workflow_execution_id") REFERENCES "workflow_executions"("id") ON DELETE CASCADE
);

CREATE INDEX ON "workflow_executions" ("workflow_id");
CREATE INDEX ON "steps" ("workflow_execution_id");
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS "steps";
DROP TABLE IF EXISTS "workflow_executions";
DROP TABLE IF EXISTS "workflows";
-- +goose StatementEnd
