-- name: StepExecutionGet :one
SELECT * FROM step_executions
WHERE id = @id;

-- name: StepExecutionListByWorkflowExecutionID :many
SELECT * FROM step_executions
WHERE workflow_execution_id = @workflow_execution_id;

-- name: StepExecutionBatchInsert :copyfrom
INSERT INTO step_executions (
	id,
	workflow_execution_id,
	status,
	node,
	inputs,
	outputs,
	error,
	created_at,
	updated_at,
	started_at,
	completed_at
)
VALUES (
	@id,
	@workflow_execution_id,
	@status,
	@node,
	@inputs,
	@outputs,
	@error,
	@created_at,
	@updated_at,
	@started_at,
	@completed_at
);

-- name: StepExecutionUpdate :one
UPDATE step_executions
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
