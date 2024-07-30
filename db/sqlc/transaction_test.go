package db

import (
	"banking_application/util"
	"context"
	"log"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func CreateRandomTransaction(t *testing.T, a Account) Transaction {
	arg := CreateTransactionParams{
		AccountNumber: a.AccountNumber,
		Amount:        util.RandomTransactionAmount(),
	}
	transaction, err := testQueries.CreateTransaction(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, transaction)

	require.NotEmpty(t, transaction.TransactionID)
	require.NotEmpty(t, transaction.CreatedAt)

	require.Equal(t, transaction.AccountNumber, arg.AccountNumber)
	require.Equal(t, transaction.Amount, arg.Amount)

	return transaction
}

func TestCreateTransaction(t *testing.T) {
	account := CreateRandomAccount(t)
	CreateRandomTransaction(t, account)
}

func TestGetTransaction(t *testing.T) {
	account1 := CreateRandomAccount(t)
	transaction1 := CreateRandomTransaction(t, account1)
	transaction2, err := testQueries.GetTransaction(context.Background(), transaction1.TransactionID)

	require.NoError(t, err)
	require.NotEmpty(t, transaction2)

	log.Printf("Transaction: %+v\n", transaction2)

	require.Equal(t, transaction1.AccountNumber, transaction2.AccountNumber)
	require.Equal(t, transaction1.TransactionID, transaction2.TransactionID)
	require.WithinDuration(t, transaction1.CreatedAt, transaction2.CreatedAt, time.Second)
	require.Equal(t, transaction1.Amount, transaction1.Amount)
}

func TestListTransactions(t *testing.T) {
	account := CreateRandomAccount(t)

	for i := 0; i < 5; i++ {
		CreateRandomTransaction(t, account)
	}

	arg := ListTransactionsParams{
		AccountNumber: account.AccountNumber,
		Limit:         5,
		Offset:        0,
	}

	transactions, err := testQueries.ListTransactions(context.Background(), arg)

	require.NoError(t, err)
	require.Len(t, transactions, 5)

	for _, transaction := range transactions {
		require.NotEmpty(t, transaction)
	}
}
