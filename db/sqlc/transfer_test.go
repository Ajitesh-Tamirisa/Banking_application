package db

import (
	"banking_application/util"
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func CreateRandomTransfer(t *testing.T, fromAccount Account, toAccount Account) Transfer {
	arg := CreateTransferParams{
		FromAccountNumber: fromAccount.AccountNumber,
		ToAccountNumber:   toAccount.AccountNumber,
		Amount:            util.RandomInt(1, 5),
	}

	transfer, err := testQueries.CreateTransfer(context.Background(), arg)

	require.NotEmpty(t, transfer)
	require.NoError(t, err)

	require.NotEmpty(t, transfer.TransferID)
	require.Equal(t, transfer.FromAccountNumber, arg.FromAccountNumber)
	require.Equal(t, transfer.ToAccountNumber, arg.ToAccountNumber)
	require.Equal(t, transfer.Amount, arg.Amount)
	require.NotEmpty(t, transfer.CreatedAt)

	return transfer
}

func TestCreateTransfer(t *testing.T) {
	fromAccount := CreateRandomAccount(t)
	toAccount := CreateRandomAccount(t)
	CreateRandomTransfer(t, fromAccount, toAccount)
}

func TestGetTransfer(t *testing.T) {
	fromAccount := CreateRandomAccount(t)
	toAccount := CreateRandomAccount(t)
	transfer1 := CreateRandomTransfer(t, fromAccount, toAccount)

	transfer2, err := testQueries.GetTransfer(context.Background(), transfer1.TransferID)

	require.NoError(t, err)
	require.NotEmpty(t, transfer2)

	require.Equal(t, transfer1.FromAccountNumber, transfer2.FromAccountNumber)
	require.Equal(t, transfer1.ToAccountNumber, transfer2.ToAccountNumber)
	require.Equal(t, transfer1.Amount, transfer2.Amount)
	require.WithinDuration(t, transfer1.CreatedAt, transfer2.CreatedAt, time.Second)
}

func TestListTransfers(t *testing.T) {
	fromAccount := CreateRandomAccount(t)
	toAccount := CreateRandomAccount(t)

	for i := 0; i < 3; i++ {
		CreateRandomTransfer(t, fromAccount, toAccount)
	}

	arg := ListTransfersParams{
		FromAccountNumber: fromAccount.AccountNumber,
		ToAccountNumber:   toAccount.AccountNumber,
		Limit:             3,
		Offset:            0,
	}

	transfers, err := testQueries.ListTransfers(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, transfers)

	for _, transfer := range transfers {
		require.NotEmpty(t, transfer)
		require.Equal(t, transfer.FromAccountNumber, fromAccount.AccountNumber)
		require.Equal(t, transfer.ToAccountNumber, toAccount.AccountNumber)
	}
}
