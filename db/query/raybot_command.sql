-- name: GetRaybotCommand :one
SELECT * FROM raybot_commands
WHERE id = $1;

-- name: CreateRaybotCommand :exec
INSERT INTO raybot_commands (
    id,
    raybot_id,
    type,
    status,
    input,
    created_at
)
VALUES (
    $1,
    $2,
    $3,
    $4,
    $5,
    $6
);

-- name: UpdateRaybotCommand :exec
UPDATE raybot_commands
SET status = $1,
    output = $2,
    completed_at = $3
WHERE id = $4;

-- name: DeleteRaybotCommand :exec
DELETE FROM raybot_commands
WHERE id = $1;
