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

-- name: GetAccountForUpdate :one
SELECT * FROM accounts
WHERE account_id=$1 LIMIT 1
FOR NO KEY UPDATE;
-- To Avoid DeadLock Happening between operations caused by foreign key

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

-- name: AddAccountBalance :one
UPDATE accounts
SET balance = balance + sqlc.arg(amount)
WHERE account_id = sqlc.arg(account_id)
RETURNING *;


-- name: DeleteAccount :exec
DELETE FROM accounts
WHERE account_id = $1;