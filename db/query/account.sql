-- name: CreateAccount :one
INSERT INTO accounts (user_name, email, balance, currency) 
VALUES ($1, $2, $3, $4) 
RETURNING *; 

-- name: GetAccount :one
SELECT * FROM accounts
WHERE account_number = $1 LIMIT 1;

-- name: GetAccountForUpdate :one
SELECT * FROM accounts
WHERE account_number = $1 LIMIT 1 
FOR NO KEY UPDATE;

-- name: ListAccounts :many
SELECT * FROM accounts
ORDER BY account_number
LIMIT $1
OFFSET $2;

-- name: UpdateAccount :one
UPDATE accounts
SET balance = $2
WHERE account_number = $1
RETURNING *;

-- name: AddAccountBalance :one
UPDATE accounts
SET balance = balance + sqlc.arg(amount)
WHERE account_number = sqlc.arg(accountNumber)
RETURNING *;

-- name: DeleteAccount :exec
DELETE FROM accounts WHERE account_number = $1;