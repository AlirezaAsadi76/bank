package db

import (
	"context"
	"database/sql"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func createFakeEntry(t *testing.T) Entry {
	acc := createFakeAccount(t)
	arg := CreateEntryParams{AccountID: acc.ID, Amount: createRandomInteger(1000, 100000)}
	ent, err := queryTest.CreateEntry(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, ent)

	require.Equal(t, arg.AccountID, ent.AccountID)
	require.Equal(t, arg.Amount, ent.Amount)

	require.NotZero(t, ent.ID)
	require.NotZero(t, ent.CreatedAt)
	return ent
}

func TestQueries_CreateEntry(t *testing.T) {
	createFakeEntry(t)
}

func TestQueries_GetEntry(t *testing.T) {
	ent := createFakeEntry(t)

	ent2, err := queryTest.GetEntry(context.Background(), ent.ID)

	require.NoError(t, err)
	require.NotEmpty(t, ent2)

	require.Equal(t, ent2.AccountID, ent.AccountID)
	require.Equal(t, ent2.Amount, ent.Amount)

	require.WithinDuration(t, ent2.CreatedAt, ent.CreatedAt, time.Second)

}

func TestQueries_UpdateEntry(t *testing.T) {
	ent := createFakeEntry(t)

	arg := UpdateEntryParams{ID: ent.ID, Amount: createRandomInteger(1000, 100000)}

	ent2, err := queryTest.UpdateEntry(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, ent2)

	require.Equal(t, ent2.AccountID, ent.AccountID)
	require.Equal(t, ent2.Amount, arg.Amount)

	require.WithinDuration(t, ent2.CreatedAt, ent.CreatedAt, time.Second)
}

func TestQueries_DeleteEntry(t *testing.T) {
	ent := createFakeEntry(t)
	err := queryTest.DeleteEntry(context.Background(), ent.ID)
	require.NoError(t, err)
	ent2, err := queryTest.GetEntry(context.Background(), ent.ID)
	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, ent2)
}

func TestQueries_ListEntries(t *testing.T) {

	for i := 0; i < 10; i++ {
		createFakeEntry(t)
	}
	arg := ListEntriesParams{Limit: 8, Offset: 1}

	ents, err := queryTest.ListEntries(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, ents)
	require.Len(t, ents, 8)

	for _, ent := range ents {
		require.NotEmpty(t, ent)
	}

}
