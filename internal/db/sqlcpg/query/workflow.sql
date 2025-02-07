-- name: WorkflowGetByID :one
SELECT * FROM workflows
WHERE id = @id;

-- name: WorkflowInsert :exec
INSERT INTO workflows (
	id,
	name,
	description,
	is_draft,
	is_valid,
	data,
	created_at,
	updated_at
)
VALUES (
	@id,
	@name,
	@description,
	@is_draft,
	@is_valid,
	@data,
	@created_at,
	@updated_at
);

-- name: WorkflowUpdate :one
UPDATE workflows
SET
	name = CASE WHEN @set_name::boolean THEN @name ELSE name END,
	description = CASE WHEN @set_description::boolean THEN @description ELSE description END,
	is_draft = CASE WHEN @set_is_draft::boolean THEN @is_draft ELSE is_draft END,
	is_valid = CASE WHEN @set_is_valid::boolean THEN @is_valid ELSE is_valid END,
	data = CASE WHEN @set_data::boolean THEN @data ELSE data END,
	updated_at = NOW()
WHERE id = @id
RETURNING *;

-- name: WorkflowDelete :exec
DELETE FROM workflows
WHERE id = @id;
