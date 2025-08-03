INSERT INTO audit_records (
	id, request_id, action, created_at, data, origin_ip, resource_id, resource_type, user_id
) VALUES (
	$1, $2, $3, $4, $5, $6, $7, $8, $9
)
