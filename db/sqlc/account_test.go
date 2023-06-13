package db

import (
	"context"
	"database/sql"
	"simplebank/util"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func CreateRandomAccount(t *testing.T) Account {
	arg := CreateAccountParams{
		Owner: util.RandomOwner(),
		Balance: util.RandomMoney(),
		Currency: util.RandomCurrency(),
	}

	account, err := testQueries.CreateAccount(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, account)

	require.Equal(t, arg.Owner, account.Owner)
	require.Equal(t, arg.Balance, account.Balance)
	require.Equal(t, arg.Currency,account.Currency)

	require.NotZero(t, account.ID)
	require.NotZero(t, account.CreatedAt)
	return account
}

func TestCreateAccount(t *testing.T) {
	CreateRandomAccount(t)	
}

func TestGetAccount(t *testing.T) {
	expected := CreateRandomAccount(t)
	actual, err := testQueries.GetAccount(context.Background(), expected.ID)
	require.NoError(t, err)
	require.NotEmpty(t, actual)

	require.Equal(t, actual.ID, expected.ID)
	require.Equal(t, actual.Owner, expected.Owner)
	require.Equal(t, actual.Balance, expected.Balance)
	require.Equal(t, actual.Currency, expected.Currency)
	require.WithinDuration(t, actual.CreatedAt, expected.CreatedAt, time.Second)
}

func TestUpdateAccount(t *testing.T) {
	account1 := CreateRandomAccount(t)
	arg := UpdateAccountParams{
		ID: account1.ID,
		Balance: account1.Balance,
	}
	
	account2, err := testQueries.UpdateAccount(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, account2)
	require.Equal(t, account2.ID, arg.ID)
	require.Equal(t, account2.Balance, arg.Balance)
	require.Equal(t, account2.Owner, account1.Owner)
	require.Equal(t, account2.Currency, account1.Currency)
	require.WithinDuration(t, account2.CreatedAt, account1.CreatedAt, time.Second)
}

func TestDeleteAccount(t *testing.T) {
	account1 := CreateRandomAccount(t)
	err := testQueries.DeleteAccount(context.Background(), account1.ID)
	require.NoError(t, err)

	account2, err := testQueries.GetAccount(context.Background(), account1.ID)
	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, account2)
}

func TestListAccount(t *testing.T) {
	for i := 0; i < 10; i++ {
		CreateRandomAccount(t)
	}

	arg := ListAccountsParams{
		Limit: 5,
		Offset: 5,
	}
	accounts, err := testQueries.ListAccounts(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, accounts, 5)

	for _, account := range accounts {
		require.NotEmpty(t, account)
	}
}