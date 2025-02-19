-- name: WorkflowExecutionGetByID :one
SELECT * FROM workflow_executions
WHERE id = @id;

-- name: WorkflowExecutionInsert :exec
INSERT INTO workflow_executions (
	id,
	workflow_id,
	status,
	data,
	inputs,
	outputs,
	error,
	created_at,
	started_at,
	completed_at
)
VALUES (
	@id,
	@workflow_id,
	@status,
	@data,
	@inputs,
	@outputs,
	@error,
	@created_at,
	@started_at,
	@completed_at
);

-- name: WorkflowExecutionUpdate :one
UPDATE workflow_executions
SET
	status = CASE WHEN @set_status::boolean THEN @status ELSE status END,
	inputs = CASE WHEN @set_inputs::boolean THEN @inputs ELSE inputs END,
	outputs = CASE WHEN @set_outputs::boolean THEN @outputs ELSE outputs END,
	error = CASE WHEN @set_error::boolean THEN @error ELSE error END,
	started_at = CASE WHEN @set_started_at::boolean THEN @started_at ELSE started_at END,
	completed_at = CASE WHEN @set_completed_at::boolean THEN @completed_at ELSE completed_at END,
	updated_at = NOW()
WHERE id = @id
RETURNING *;
