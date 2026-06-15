-- name: Validate :one
SELECT *
FROM sessions
LEFT JOIN users ON sessions.user_id = users.user_id
WHERE session_id = @sessionId;
