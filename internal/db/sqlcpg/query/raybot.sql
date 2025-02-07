-- name: RaybotGetByID :one
SELECT * FROM raybots
WHERE id = @id;

-- name: RaybotInsert :exec
INSERT INTO raybots (
    id,
    name,
	ip_address,
	last_connected_at,
	is_online,
	control_mode,
    created_at,
    updated_at
)
VALUES (
    @id,
    @name,
    @ip_address,
    @last_connected_at,
    @is_online,
    @control_mode,
    @created_at,
    @updated_at
);

-- name: RaybotUpdate :one
UPDATE raybots
SET
	name = CASE WHEN @set_name::boolean THEN @name ELSE name END,
	ip_address = CASE WHEN @set_ip_address::boolean THEN @ip_address ELSE ip_address END,
	last_connected_at = CASE WHEN @set_last_connected_at::boolean THEN @last_connected_at ELSE last_connected_at END,
	is_online = CASE WHEN @set_is_online::boolean THEN @is_online ELSE is_online END,
	control_mode = CASE WHEN @set_control_mode::boolean THEN @control_mode ELSE control_mode END,
	updated_at = now()
WHERE id = @id
RETURNING *;

-- name: RaybotDelete :exec
DELETE FROM raybots
WHERE id = @id;
