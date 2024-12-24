-- name: GetQRLocation :one
SELECT * FROM qr_locations
WHERE id = $1;

-- name: ExistsQRLocationByQRCode :one
SELECT EXISTS(
	SELECT 1
	FROM qr_locations
	WHERE qr_code = $1
);

-- name: CreateQRLocation :exec
INSERT INTO qr_locations (
    id,
    name,
    qr_code,
	metadata,
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

-- name: UpdateQRLocation :one
UPDATE qr_locations
SET name = $1,
    qr_code = $2,
	metadata = $3,
    updated_at = $4
WHERE id = $5
RETURNING *;

-- name: DeleteQRLocation :exec
DELETE FROM qr_locations
WHERE id = $1;
