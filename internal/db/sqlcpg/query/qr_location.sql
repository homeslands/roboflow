-- name: QRLocationGetByID :one
SELECT * FROM qr_locations
WHERE id = @id;

-- name: QRLocationInsert :exec
INSERT INTO qr_locations (
    id,
    name,
    qr_code,
	metadata,
    created_at,
    updated_at
)
VALUES (
    @id,
    @name,
    @qr_code,
    @metadata,
    @created_at,
    @updated_at
);

-- name: QRLocationUpdate :one
UPDATE qr_locations
SET
	name = CASE WHEN @set_name::boolean THEN @name ELSE name END,
    qr_code = CASE WHEN @set_qr_code::boolean THEN @qr_code ELSE qr_code END,
	metadata = CASE WHEN @set_metadata::boolean THEN @metadata ELSE metadata END,
    updated_at = NOW()
WHERE id = @id
RETURNING *;

-- name: QRLocationDelete :exec
DELETE FROM qr_locations
WHERE id = @id;
