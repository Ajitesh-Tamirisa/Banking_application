-- name: CreateTransaction :one
INSERT INTO transactions (account_number, amount)
VALUES ($1, $2)
RETURNING *;

-- name: GetTransaction :one
SELECT * FROM transactions
WHERE transaction_id = $1 LIMIT 1;

-- name: ListTransactions :many
SELECT * FROM transactions
WHERE account_number = $1
ORDER BY transaction_id
LIMIT $2
OFFSET $3;