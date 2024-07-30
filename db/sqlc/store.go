package db

import (
	"context"
	"database/sql"
	"fmt"
	"log"
)

type Store struct {
	*Queries
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{
		db:      db,
		Queries: New(db),
	}
}

func (store *Store) execTx(ctx context.Context, fn func(*Queries) error) error {
	//Executes a function within a DB transaction

	tx, err := store.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	q := New(tx)
	err = fn(q)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("tx err: %v, rb err: %v", err, rbErr)
		}
		return err
	}
	log.Println("Transaction committed successfully")
	return tx.Commit()
}

type TransferTxParams struct {
	FromAccountNumber int64 `json:"from_account_number"`
	ToAccountNumber   int64 `json:"to_account_number"`
	Amount            int64 `json:"amount"`
}

type TransferTxResult struct {
	Transfer        Transfer    `json:"transfer"`
	FromAccount     Account     `json:"from_account"`
	ToAccount       Account     `json:"to_account"`
	FromTransaction Transaction `json:"from_transaction"`
	ToTransaction   Transaction `json:"to_transaction"`
}

func (store *Store) TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error) {
	//Creates a transfer record, add account transaction record, and update accounts' balance within a single database transaction
	var result TransferTxResult

	err := store.execTx(ctx, func(q *Queries) error {
		var err error

		result.Transfer, err = q.CreateTransfer(ctx, CreateTransferParams{
			FromAccountNumber: arg.FromAccountNumber,
			ToAccountNumber:   arg.ToAccountNumber,
			Amount:            arg.Amount,
		})
		if err != nil {
			return err
		}

		result.FromTransaction, err = q.CreateTransaction(ctx, CreateTransactionParams{
			AccountNumber: arg.FromAccountNumber,
			Amount:        -arg.Amount,
		})
		if err != nil {
			return err
		}

		result.ToTransaction, err = q.CreateTransaction(ctx, CreateTransactionParams{
			AccountNumber: arg.ToAccountNumber,
			Amount:        arg.Amount,
		})
		if err != nil {
			return err
		}

		if arg.FromAccountNumber < arg.ToAccountNumber {
			result.FromAccount, result.ToAccount, err = addMoney(ctx, q, arg.FromAccountNumber, -arg.Amount, arg.ToAccountNumber, arg.Amount)
		} else {
			result.ToAccount, result.FromAccount, err = addMoney(ctx, q, arg.ToAccountNumber, arg.Amount, arg.FromAccountNumber, -arg.Amount)
		}
		if err != nil {
			return err
		}

		return nil
	})
	return result, err
}

func addMoney(
	ctx context.Context,
	q *Queries,
	accountNumber1 int64,
	amount1 int64,
	accountNumber2 int64,
	amount2 int64) (account1 Account, account2 Account, err error) {
	account1, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
		Accountnumber: accountNumber1,
		Amount:        amount1,
	})
	if err != nil {
		return
	}

	account2, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
		Accountnumber: accountNumber2,
		Amount:        amount2,
	})
	if err != nil {
		return
	}
	return
}
