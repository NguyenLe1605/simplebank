package db

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTransferTx(t *testing.T) {
	store := NewStore(testDB)
	fromAccount := CreateRandomAccount(t)
	toAccount := CreateRandomAccount(t)

	fmt.Println(">>before: ", fromAccount.Balance, toAccount.Balance);

	n := 3
	amount := int64(10)

	params := TransferTxParams{
		ToAccountID: toAccount.ID,
		FromAccountID: fromAccount.ID,
		Amount: amount,
	}

	err := make(chan error)
	result := make(chan TransferTxResult)
	for i := 0 ; i < n; i++ {
		go func() {
			ctx := context.Background()
			res, er := store.TransferTx(ctx, params)
			err <- er
			result <- res
		}()
	}

	existed := make(map[int64]bool)

	for i := 0; i < n; i++ {
		er := <- err
		require.NoError(t, er)
		res := <- result
		require.NotEmpty(t, res)

		// check transfer
		transfer := res.Transfer
		require.NotEmpty(t, transfer)
		require.Equal(t, toAccount.ID, res.Transfer.ToAccountID)
		require.Equal(t, fromAccount.ID, res.Transfer.FromAccountID)
		require.Equal(t, params.Amount, res.Transfer.Amount)
		require.NotZero(t, transfer.ID)
		require.NotZero(t, transfer.CreatedAt)
		_, er = store.GetTransfer(context.Background(), transfer.ID)
		require.NoError(t, er)
		
		// check entries
		fromEntry := res.FromEntry
		require.Equal(t, fromAccount.ID, res.FromEntry.AccountID)
		require.Equal(t, params.Amount, -res.FromEntry.Amount)
		require.NotZero(t, fromEntry.ID)
		require.NotZero(t, fromEntry.CreatedAt)
		_, er = store.GetEntry(context.Background(), fromEntry.ID)
		require.NoError(t, er)

		toEntry := res.ToEntry
		require.Equal(t, toAccount.ID, res.ToEntry.AccountID)
		require.Equal(t, params.Amount, res.ToEntry.Amount)
		require.NotZero(t, toEntry.ID)
		require.NotZero(t, toEntry.CreatedAt)
		_, er = store.GetEntry(context.Background(), toEntry.ID)
		require.NoError(t, er)

		// Check accounts
		fromTestAcc := res.FromAccount
		require.NotEmpty(t,fromTestAcc)
		require.Equal(t, fromTestAcc.ID, fromAccount.ID)

		toTestAcc := res.ToAccount
		require.NotEmpty(t,toTestAcc)
		require.Equal(t, toTestAcc.ID, toAccount.ID)

		fmt.Println(">>>tx: ", fromTestAcc.Balance, toTestAcc.Balance)

		// Check balance
		diff1 := fromAccount.Balance - fromTestAcc.Balance
		diff2 := toTestAcc.Balance - toAccount.Balance
		require.Equal(t, diff1, diff2)
		require.True(t, diff1 > 0)
		require.True(t, diff1 % amount == 0)

		k := diff1 / amount
		require.True(t, k >= 1 && k <= int64(n))
		require.NotContains(t, existed, k)
		existed[k] = true
	}

	updatedAccount1, e := store.GetAccount(context.Background(), fromAccount.ID)
	require.NoError(t, e)

	updatedAccount2, e := store.GetAccount(context.Background(), toAccount.ID)
	require.NoError(t, e)

	fmt.Println(">>>after: ", updatedAccount1.Balance, updatedAccount2.Balance)

	require.Equal(t, fromAccount.Balance - int64(n) * amount, updatedAccount1.Balance)
	require.Equal(t, toAccount.Balance + int64(n) * amount, updatedAccount2.Balance)
}

func TestTransferTxDeadlock(t *testing.T) {
	store := NewStore(testDB)
	account1 := CreateRandomAccount(t)
	account2 := CreateRandomAccount(t)

	fmt.Println(">>before: ", account1.Balance, account2.Balance);

	n := 10
	amount := int64(10)


	err := make(chan error)
	for i := 0 ; i < n; i++ {
		fromAccountID, toAccountID := account1.ID, account2.ID;
		if (i % 2 == 1) {
			toAccountID, fromAccountID = account1.ID, account2.ID;	
		}
		go func() {
			ctx := context.Background()
			_, er := store.TransferTx(ctx, TransferTxParams{
				FromAccountID: fromAccountID,
				ToAccountID: toAccountID,
				Amount: amount,
			})
			err <- er
		}()
	}


	for i := 0; i < n; i++ {
		er := <- err
		require.NoError(t, er)
	}

	updatedAccount1, e := store.GetAccount(context.Background(), account1.ID)
	require.NoError(t, e)

	updatedAccount2, e := store.GetAccount(context.Background(), account2.ID)
	require.NoError(t, e)

	fmt.Println(">>>after: ", updatedAccount1.Balance, updatedAccount2.Balance)

	require.Equal(t, account1.Balance, updatedAccount1.Balance)
	require.Equal(t, account2.Balance, updatedAccount2.Balance)
}

