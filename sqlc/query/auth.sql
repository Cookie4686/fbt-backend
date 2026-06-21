-- name: Validate :one
SELECT *
FROM sessions
LEFT JOIN users ON sessions.user_id = users.user_id
WHERE session_id = @sessionId;

-- name: GetAll :many
SELECT transactions.transaction_id, transactions.datetime, entries.account_id, entries.amount
FROM accounts
JOIN entries ON entries.account_id = accounts.account_id
JOIN transactions ON transactions.transaction_id = entries.transaction_id
WHERE accounts.user_id = @user_id;