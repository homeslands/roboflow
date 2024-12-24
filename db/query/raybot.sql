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

-- name: UpdateRaybot :exec
UPDATE raybots
SET name = $1,
    token = $2,
    status = $3,
    updated_at = $4
WHERE id = $5;

-- name: UpdateRaybotStatus :exec
UPDATE raybots
SET status = $1,
    updated_at = $2
WHERE id = $3;

-- name: DeleteRaybot :exec
DELETE FROM raybots
WHERE id = $1;
