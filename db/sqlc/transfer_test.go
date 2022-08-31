package db

import (
	"context"
	"database/sql"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func createFakeTransfer(t *testing.T) Transfer {
	acc1 := createFakeAccount(t)
	acc2 := createFakeAccount(t)

	arg := CreateTransferParams{FromAccountID: acc1.ID, ToAccountID: acc2.ID, Amount: createRandomInteger(1000, 100000)}

	tra, err := queryTest.CreateTransfer(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, tra)
	require.Equal(t, tra.FromAccountID, arg.FromAccountID)
	require.Equal(t, tra.ToAccountID, arg.ToAccountID)
	require.Equal(t, tra.Amount, arg.Amount)

	require.NotZero(t, tra.ID)
	require.NotZero(t, tra.CreatedAt)

	return tra
}
func TestQueries_CreateTransfer(t *testing.T) {

	createFakeTransfer(t)
}

func TestQueries_GetTransfer(t *testing.T) {
	tra := createFakeTransfer(t)

	tra2, err := queryTest.GetTransfer(context.Background(), tra.ID)

	require.NoError(t, err)
	require.NotEmpty(t, tra2)

	require.Equal(t, tra.FromAccountID, tra2.FromAccountID)
	require.Equal(t, tra.ToAccountID, tra2.ToAccountID)
	require.Equal(t, tra.Amount, tra2.Amount)

	require.WithinDuration(t, tra.CreatedAt, tra2.CreatedAt, time.Second)
}

func TestQueries_UpdateTransfer(t *testing.T) {
	tra := createFakeTransfer(t)

	arg := UpdateTransferParams{ID: tra.ID, Amount: createRandomInteger(1000, 10000)}

	tra2, err := queryTest.UpdateTransfer(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, tra2)

	require.Equal(t, tra.FromAccountID, tra2.FromAccountID)
	require.Equal(t, tra.ToAccountID, tra2.ToAccountID)
	require.Equal(t, arg.Amount, tra2.Amount)

	require.WithinDuration(t, tra.CreatedAt, tra2.CreatedAt, time.Second)

}

func TestQueries_DeleteTransfer(t *testing.T) {
	tra := createFakeTransfer(t)

	err := queryTest.DeleteTransfer(context.Background(), tra.ID)

	require.NoError(t, err)

	tra2, err := queryTest.GetTransfer(context.Background(), tra.ID)

	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, tra2)

}

func TestQueries_ListTransfers(t *testing.T) {
	for i := 0; i < 10; i++ {
		createFakeTransfer(t)
	}

	arg := ListTransfersParams{Limit: 2, Offset: 5}

	tras, err := queryTest.ListTransfers(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, tras)
	require.Len(t, tras, 2)

	for _, tra := range tras {
		require.NotEmpty(t, tra)
	}
}
