package db

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTransferTx(t *testing.T) {
	store := NewStore(testDB)

	account1 := CreateRandomAccount(t)
	account2 := CreateRandomAccount(t)

	//run n concurrent transfer transactions
	n := 5
	amount := int64(10)

	errs := make(chan error)
	results := make(chan TransferTxResult)

	for i := 0; i < n; i++ {
		go func() {
			ctx := context.Background()
			result, err := store.TransferTx(ctx, TransferTxParams{
				FromAccountNumber: account1.AccountNumber,
				ToAccountNumber:   account2.AccountNumber,
				Amount:            amount,
			})
			errs <- err
			results <- result
		}()
	}

	//check results
	existed := make(map[int]bool)
	for i := 0; i < n; i++ {
		err := <-errs
		require.NoError(t, err)

		result := <-results
		require.NotEmpty(t, result)

		//check transfer
		transfer := result.Transfer
		require.NotEmpty(t, transfer)
		require.Equal(t, account1.AccountNumber, transfer.FromAccountNumber)
		require.Equal(t, account2.AccountNumber, transfer.ToAccountNumber)
		require.Equal(t, amount, transfer.Amount)
		require.NotZero(t, transfer.TransferID)
		require.NotZero(t, transfer.CreatedAt)

		_, err = store.GetTransfer(context.Background(), transfer.TransferID)
		require.NoError(t, err)

		//check transactionss
		fromTransaction := result.FromTransaction
		require.NotEmpty(t, fromTransaction)
		require.Equal(t, account1.AccountNumber, fromTransaction.AccountNumber)
		require.Equal(t, -amount, fromTransaction.Amount)
		require.NotZero(t, fromTransaction.AccountNumber)
		require.NotZero(t, fromTransaction.CreatedAt)

		_, err = store.GetTransaction(context.Background(), fromTransaction.TransactionID)
		require.NoError(t, err)

		toTransaction := result.ToTransaction
		require.NotEmpty(t, toTransaction)
		require.Equal(t, account2.AccountNumber, toTransaction.AccountNumber)
		require.Equal(t, amount, toTransaction.Amount)
		require.NotZero(t, toTransaction.AccountNumber)
		require.NotZero(t, toTransaction.CreatedAt)

		_, err = store.GetTransaction(context.Background(), toTransaction.TransactionID)
		require.NoError(t, err)

		//TODO: Check account balance

		//check accounts
		fromAccount := result.FromAccount
		require.NotEmpty(t, fromAccount)
		require.Equal(t, account1.AccountNumber, fromAccount.AccountNumber)

		toAccount := result.ToAccount
		require.NotEmpty(t, toAccount)
		require.Equal(t, account2.AccountNumber, toAccount.AccountNumber)

		//check accounts' balances
		diff1 := account1.Balance - fromAccount.Balance
		diff2 := toAccount.Balance - account2.Balance

		require.Equal(t, diff1, diff2)
		require.True(t, diff1 > 0)
		require.True(t, diff1%amount == 0)

		k := int(diff1 / amount)
		require.True(t, k >= 1 && k <= n)
		require.NotContains(t, existed, k)
		existed[k] = true
	}

	//check the final updated balances
	updateAccount1, err := testQueries.GetAccount(context.Background(), account1.AccountNumber)
	require.NoError(t, err)

	updateAccount2, err := testQueries.GetAccount(context.Background(), account2.AccountNumber)
	require.NoError(t, err)

	require.Equal(t, account1.Balance-int64(n)*amount, updateAccount1.Balance)
	require.Equal(t, account2.Balance+int64(n)*amount, updateAccount2.Balance)
}

func TestTransferTxDeadlock(t *testing.T) {
	store := NewStore(testDB)

	account1 := CreateRandomAccount(t)
	account2 := CreateRandomAccount(t)

	//run n concurrent transfer transactions
	n := 10
	amount := int64(10)

	errs := make(chan error)

	for i := 0; i < n; i++ {
		fromAccountNumber := account1.AccountNumber
		toAccountNumber := account2.AccountNumber

		if i%2 == 1 {
			fromAccountNumber = account2.AccountNumber
			toAccountNumber = account1.AccountNumber
		}

		go func() {
			ctx := context.Background()
			_, err := store.TransferTx(ctx, TransferTxParams{
				FromAccountNumber: fromAccountNumber,
				ToAccountNumber:   toAccountNumber,
				Amount:            amount,
			})
			errs <- err
		}()
	}

	//check results
	for i := 0; i < n; i++ {
		err := <-errs
		require.NoError(t, err)
	}

	//check the final updated balances
	updateAccount1, err := testQueries.GetAccount(context.Background(), account1.AccountNumber)
	require.NoError(t, err)

	updateAccount2, err := testQueries.GetAccount(context.Background(), account2.AccountNumber)
	require.NoError(t, err)

	require.Equal(t, account1.Balance, updateAccount1.Balance)
	require.Equal(t, account2.Balance, updateAccount2.Balance)
}
