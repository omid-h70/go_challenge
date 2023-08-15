-- name: CreateAccount :one
INSERT INTO accounts(
owner,
balance,
currency
) VALUES (
$1, $2, $3
)RETURNING *;
-- name: ListAccounts :many
SELECT * FROM accounts
ORDER BY id
LIMIT $1
OFFSET $2;