package db

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestTransferTx(t *testing.T) {
	store := NewStore(testDB)
	account1 := CreateRandomAccount(t)
	account2 := CreateRandomAccount(t)
	fmt.Println(">> before:", account1.Balance, account2.Balance)
	// run concurrent to do transfer
	n := 5
	amount := int64(10)
	errors := make(chan error)
	results := make(chan TransferTxResult)
	for i := 0; i < n; i++ {
		txName := fmt.Sprintf("transfer-%d", i+1)
		go func() {
			ctx := context.WithValue(context.Background(), txKey, txName)
			res, err := store.TransferTx(ctx, TransferTxParams{
				FromAccountID: account1.ID,
				ToAccountID:   account2.ID,
				Amount:        amount,
			})
			errors <- err
			results <- res
		}()
	}

	existed := make(map[int]bool)
	for i := 0; i < n; i++ {
		err := <-errors
		require.NoError(t, err)
		res := <-results

		transfer := res.Transfer
		require.NotEmpty(t, res)
		require.Equal(t, transfer.FromAccountID, account1.ID)
		require.Equal(t, transfer.ToAccountID, account2.ID)
		require.NotEmpty(t, res.Transfer)
		require.NotZero(t, res.Transfer.CreatedAt)
		require.NotZero(t, res.Transfer.ID)
		_, err = store.GetTransfer(context.Background(), res.Transfer.ID)
		require.NoError(t, err)

		// check entry
		require.NotEmpty(t, res.FromEntry)
		require.Equal(t, res.FromEntry.AccountID, account1.ID)
		require.Equal(t, res.FromEntry.Amount, -amount)
		require.NotZero(t, res.FromEntry.CreatedAt)

		_, err = store.GetEntry(context.Background(), res.FromEntry.ID)
		require.NoError(t, err)

		require.NotEmpty(t, res.ToEntry)
		require.Equal(t, res.ToEntry.AccountID, account2.ID)
		require.Equal(t, res.ToEntry.Amount, amount)
		require.NotZero(t, res.ToEntry.CreatedAt)

		_, err = store.GetEntry(context.Background(), res.ToEntry.ID)
		require.NoError(t, err)

		//check account
		fromAccount := res.FromAccount
		require.NotEmpty(t, fromAccount)
		require.Equal(t, fromAccount.ID, account1.ID)

		toAccount := res.ToAccount
		require.NotEmpty(t, toAccount)
		require.Equal(t, toAccount.ID, account2.ID)
		fmt.Println("tx>>", fromAccount.Balance, toAccount.Balance)

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

	//// check balance
	account1Finished, err := store.GetAccount(context.Background(), account1.ID)
	require.NoError(t, err)
	account2Finished, err := store.GetAccount(context.Background(), account2.ID)
	require.NoError(t, err)

	require.Equal(t, account1Finished.Balance, account1.Balance-int64(n)*amount)
	require.Equal(t, account2Finished.Balance, account2.Balance+int64(n)*amount)

}

func TestTransferTxDeadLock(t *testing.T) {
	store := NewStore(testDB)
	account1 := CreateRandomAccount(t)
	account2 := CreateRandomAccount(t)
	fmt.Println(">> before:", account1.Balance, account2.Balance)
	// run concurrent to do transfer
	n := 10
	amount := int64(10)
	errors := make(chan error)
	for i := 0; i < n; i++ {
		txName := fmt.Sprintf("transfer-%d", i+1)
		fromAccountId := account1.ID
		toAccountId := account2.ID

		if i%2 == 1 {
			fromAccountId = account2.ID
			toAccountId = account1.ID
		}
		go func() {
			ctx := context.WithValue(context.Background(), txKey, txName)
			_, err := store.TransferTx(ctx, TransferTxParams{
				FromAccountID: fromAccountId,
				ToAccountID:   toAccountId,
				Amount:        amount,
			})
			errors <- err
		}()
	}

	for i := 0; i < n; i++ {
		err := <-errors
		require.NoError(t, err)
	}

	//// check balance
	account1Finished, err := store.GetAccount(context.Background(), account1.ID)
	require.NoError(t, err)
	account2Finished, err := store.GetAccount(context.Background(), account2.ID)
	require.NoError(t, err)

	require.Equal(t, account1Finished.Balance, account1.Balance)
	require.Equal(t, account2Finished.Balance, account2.Balance)

}
