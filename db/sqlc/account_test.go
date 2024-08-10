package db

import (
	"banking_application/util"
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func CreateRandomAccount(t *testing.T) Account {

	user := createRandomUser(t)

	arg := CreateAccountParams{
		UserName: user.UserName,
		Email:    user.Email,
		Balance:  util.RandomMoney(),
		Currency: util.RandomCurrency(),
	}

	account, err := testQueries.CreateAccount(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, account)

	require.Equal(t, arg.UserName, account.UserName)
	require.Equal(t, arg.Email, account.Email)
	require.Equal(t, arg.Balance, account.Balance)
	require.Equal(t, arg.Currency, account.Currency)

	require.NotZero(t, account.AccountNumber)
	require.NotZero(t, account.CreatedAt)

	return account

}

func TestCreateAccount(t *testing.T) {
	CreateRandomAccount(t)
}

func TestGetAccount(t *testing.T) {
	account1 := CreateRandomAccount(t)
	account2, err := testQueries.GetAccount(context.Background(), account1.AccountNumber)

	require.NoError(t, err)
	require.NotEmpty(t, account2)

	require.Equal(t, account1.AccountNumber, account2.AccountNumber)
	require.Equal(t, account1.Balance, account2.Balance)
	require.WithinDuration(t, account1.CreatedAt, account2.CreatedAt, time.Second)
	require.Equal(t, account1.Currency, account2.Currency)
	require.Equal(t, account1.Email, account2.Email)
}

func TestUpdateAccount(t *testing.T) {
	account1 := CreateRandomAccount(t)
	arg := UpdateAccountParams{
		AccountNumber: account1.AccountNumber,
		Balance:       util.RandomMoney(),
	}

	account2, err := testQueries.UpdateAccount(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, account2)

	require.Equal(t, account1.AccountNumber, account2.AccountNumber)
	require.Equal(t, arg.Balance, account2.Balance)
	require.WithinDuration(t, account1.CreatedAt, account2.CreatedAt, time.Second)
	require.Equal(t, account1.Currency, account2.Currency)
	require.Equal(t, account1.Email, account2.Email)

}

func TestDeleteAccount(t *testing.T) {
	account1 := CreateRandomAccount(t)
	err := testQueries.DeleteAccount(context.Background(), account1.AccountNumber)

	require.NoError(t, err)
	account2, err := testQueries.GetAccount(context.Background(), account1.AccountNumber)
	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, account2)
}

func TestListAccounts(t *testing.T) {
	var lastAccount Account
	for i := 0; i < 10; i++ {
		lastAccount = CreateRandomAccount(t)
	}

	arg := ListAccountsParams{
		UserName: lastAccount.UserName,
		Limit:    5,
		Offset:   0,
	}

	accounts, err := testQueries.ListAccounts(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, accounts)

	for _, account := range accounts {
		require.NotEmpty(t, account)
		require.Equal(t, lastAccount.UserName, account.UserName)
	}
}
