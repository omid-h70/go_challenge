-- name: CreateSession :one
INSERT INTO sessions (
    session_uuid,
    user_name,
    user_agent,
    refresh_token,
    client_ip,
    is_blocked,
    expires_at
) VALUES ( $1, $2, $3, $4, $5, $6, $7 ) RETURNING *;

-- name: GetSession :one
SELECT * FROM sessions WHERE session_uuid = $1 LIMIT 1;

