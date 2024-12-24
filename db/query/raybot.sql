-- name: GetRaybot :one
SELECT * FROM raybots
WHERE id = $1;

-- name: GetRaybotStatus :one
SELECT status FROM raybots
WHERE id = $1;

-- name: ListRaybots :many
SELECT * FROM raybots
ORDER BY created_at DESC;

-- name: CreateRaybot :exec
INSERT INTO raybots (
    id,
    name,
    token,
    status,
    created_at,
    updated_at
)
VALUES (
    $1,
    $2,
    $3,
    $4,
    $5,
    $6
);

-- name: UpdateRaybot :one
UPDATE raybots
SET name = $1,
    token = $2,
    status = $3,
    ip_address = $4,
    last_connected_at = $5,
    updated_at = $6
WHERE id = $7
RETURNING *;

-- name: UpdateRaybotStatus :exec
UPDATE raybots
SET status = $1,
    updated_at = $2
WHERE id = $3;

-- name: DeleteRaybot :exec
DELETE FROM raybots
WHERE id = $1;
