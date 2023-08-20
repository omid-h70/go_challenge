-- name: CreateAccount :one
INSERT INTO accounts(
owner,
balance,
currency
) VALUES (
$1, $2, $3
)RETURNING *;

-- name: GetAccount :one
SELECT * FROM accounts
WHERE account_id=$1 LIMIT 1;

-- name: ListAccounts :many
SELECT * FROM accounts
ORDER BY account_id
LIMIT $1
OFFSET $2;

-- name: UpdateAccount :one
UPDATE accounts
SET balance = $1
WHERE account_id = $2
RETURNING *;

-- name: DeleteAccount :exec
DELETE FROM accounts
WHERE account_id = $1;