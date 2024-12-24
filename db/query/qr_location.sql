-- name: GetQRLocation :one
SELECT * FROM qr_locations
WHERE id = $1;

-- name: ListQRLocations :many
SELECT * FROM qr_locations
ORDER BY created_at DESC;

-- name: CreateQRLocation :exec
INSERT INTO qr_locations (
    id,
    name,
    qr_code,
    created_at,
    updated_at
)
VALUES (
    $1,
    $2,
    $3,
    $4,
    $5
);

-- name: UpdateQRLocation :exec
UPDATE qr_locations
SET name = $1,
    qr_code = $2,
    updated_at = $3
WHERE id = $4;

-- name: DeleteQRLocation :exec
DELETE FROM qr_locations
WHERE id = $1;
