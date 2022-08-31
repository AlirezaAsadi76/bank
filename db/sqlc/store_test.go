package db

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestStore_TransferTx(t *testing.T) {
	store := NewStore(GlobalDB)

	acc1 := createFakeAccount(t)
	acc2 := createFakeAccount(t)
	fmt.Println("Before acc1:", acc1.Balance, " acc2:", acc2.Balance)
	counter := 5
	amount := int64(100)
	errors := make(chan error)
	results := make(chan TransferTxResult)
	for i := 0; i < counter; i++ {
		go func() {
			trr, err := store.TransferTx(context.Background(), TransferTxParams{
				FromAccountID: acc1.ID,
				ToAccountID:   acc2.ID,
				Amount:        amount,
			})
			errors <- err
			results <- trr
		}()
	}

	for i := 0; i < counter; i++ {
		err := <-errors
		require.NoError(t, err)
		result := <-results
		require.NotEmpty(t, result)

		transfer := result.Transfer
		require.NotEmpty(t, transfer)
		require.Equal(t, transfer.FromAccountID, acc1.ID)
		require.Equal(t, transfer.ToAccountID, acc2.ID)
		require.Equal(t, transfer.Amount, amount)
		require.NotZero(t, transfer.ID)
		require.NotZero(t, transfer.CreatedAt)

		_, err = store.GetTransfer(context.Background(), transfer.ID)
		require.NoError(t, err)

		fromEntry := result.FromEntry
		require.NotEmpty(t, fromEntry)
		require.Equal(t, fromEntry.AccountID, acc1.ID)
		require.Equal(t, fromEntry.Amount, -amount)
		require.NotZero(t, fromEntry.ID)
		require.NotZero(t, fromEntry.CreatedAt)

		_, err = store.GetEntry(context.Background(), fromEntry.ID)
		require.NoError(t, err)

		toEntry := result.ToEntry
		require.NotEmpty(t, toEntry)
		require.Equal(t, toEntry.AccountID, acc2.ID)
		require.Equal(t, toEntry.Amount, amount)
		require.NotZero(t, toEntry.ID)
		require.NotZero(t, toEntry.CreatedAt)

		_, err = store.GetEntry(context.Background(), toEntry.ID)
		require.NoError(t, err)

		fromAccount := result.FromAccount
		require.NotEmpty(t, fromAccount)
		require.Equal(t, fromAccount.ID, acc1.ID)

		toAccount := result.ToAccount
		require.NotEmpty(t, toAccount)
		require.Equal(t, toAccount.ID, acc2.ID)

		dif1 := acc1.Balance - fromAccount.Balance
		dif2 := toAccount.Balance - acc2.Balance

		require.Equal(t, dif2, dif1)
		require.True(t, dif1 > 0)
		require.True(t, dif1%amount == 0)

		k := int(dif1 / amount)
		require.True(t, k >= 1 && k <= counter)

	}
	updateAcc1, err := queryTest.GetAccount(context.Background(), acc1.ID)
	require.NoError(t, err)
	updateAcc2, err := queryTest.GetAccount(context.Background(), acc2.ID)
	require.NoError(t, err)

	require.Equal(t, acc1.Balance-(int64(counter)*amount), updateAcc1.Balance)
	require.Equal(t, acc2.Balance+(int64(counter)*amount), updateAcc2.Balance)
	fmt.Println("After acc1:", updateAcc1.Balance, " acc2:", updateAcc2.Balance)

}

func TestStore_TransferTxDeadlock(t *testing.T) {
	store := NewStore(GlobalDB)

	acc1 := createFakeAccount(t)
	acc2 := createFakeAccount(t)
	fmt.Println("Before acc1:", acc1.Balance, " acc2:", acc2.Balance)
	counter := 10
	amount := int64(100)
	errors := make(chan error)

	for i := 0; i < counter; i++ {
		fromAcID := acc1.ID
		toAcID := acc2.ID
		if i%2 == 0 {
			fromAcID = acc2.ID
			toAcID = acc1.ID
		}
		go func() {
			_, err := store.TransferTx(context.Background(), TransferTxParams{
				FromAccountID: fromAcID,
				ToAccountID:   toAcID,
				Amount:        amount,
			})
			errors <- err

		}()
	}

	for i := 0; i < counter; i++ {
		err := <-errors
		require.NoError(t, err)

	}
	updateAcc1, err := queryTest.GetAccount(context.Background(), acc1.ID)
	require.NoError(t, err)
	updateAcc2, err := queryTest.GetAccount(context.Background(), acc2.ID)
	require.NoError(t, err)

	require.Equal(t, acc1.Balance, updateAcc1.Balance)
	require.Equal(t, acc2.Balance, updateAcc2.Balance)
	fmt.Println("After acc1:", updateAcc1.Balance, " acc2:", updateAcc2.Balance)

}
