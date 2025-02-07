-- +goose Up
-- +goose StatementBegin
CREATE TABLE "workflow_executions" (
    "id" UUID NOT NULL PRIMARY KEY,
    "workflow_id" UUID NOT NULL,
    "status" TEXT NOT NULL,
    "data" JSON NOT NULL DEFAULT '{}',
	"inputs" JSON NOT NULL DEFAULT '{}',
	"outputs" JSON NOT NULL DEFAULT '{}',
	"error" TEXT,
    "created_at" TIMESTAMPTZ NOT NULL DEFAULT NOW(),
	"updated_at" TIMESTAMPTZ NOT NULL DEFAULT NOW(),
	"started_at" TIMESTAMPTZ,
    "completed_at" TIMESTAMPTZ,

	FOREIGN KEY("workflow_id") REFERENCES "workflows"("id") ON DELETE CASCADE
);

CREATE TABLE "step_executions" (
    "id" UUID NOT NULL PRIMARY KEY,
    "workflow_execution_id" UUID NOT NULL,
    "status" TEXT NOT NULL,
    "node" JSON NOT NULL DEFAULT '{}',
	"inputs" JSON NOT NULL DEFAULT '{}',
	"outputs" JSON NOT NULL DEFAULT '{}',
	"error" TEXT,
	"created_at" TIMESTAMPTZ NOT NULL DEFAULT NOW(),
	"updated_at" TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    "started_at" TIMESTAMPTZ,
    "completed_at" TIMESTAMPTZ,

	FOREIGN KEY("workflow_execution_id") REFERENCES "workflow_executions"("id") ON DELETE CASCADE
);

CREATE INDEX ON "workflow_executions" ("workflow_id");
CREATE INDEX ON "step_executions" ("workflow_execution_id");

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS "step_executions";
DROP TABLE IF EXISTS "workflow_executions";
-- +goose StatementEnd
