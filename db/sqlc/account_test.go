package db

import (
	"context"
	"database/sql"
	"github.com/bxcodec/faker/v4"
	"github.com/stretchr/testify/require"
	"math/rand"
	"testing"
	"time"
)

func createRandomInteger(min, max int) int64 {
	return int64(rand.Intn(max-min) + min)
}
func createFakeAccount(t *testing.T) Account {
	arg := CreateAccountParams{Owner: faker.Name(),
		Balance:  createRandomInteger(1000, 100000),
		Currency: faker.Currency()}
	acc, er := queryTest.CreateAccount(context.Background(), arg)
	require.NoError(t, er)
	require.NotEmpty(t, acc)
	require.Equal(t, acc.Owner, arg.Owner)
	require.Equal(t, acc.Balance, arg.Balance)
	require.Equal(t, acc.Currency, arg.Currency)

	require.NotZero(t, acc.ID)
	require.NotZero(t, acc.CreatedAt)
	return acc
}
func TestQueries_CreateAccount(t *testing.T) {
	createFakeAccount(t)
}
func TestQueries_GetAccount(t *testing.T) {
	randAccount := createFakeAccount(t)
	acc, err := queryTest.GetAccount(context.Background(), randAccount.ID)

	require.NoError(t, err)
	require.NotEmpty(t, acc)

	require.Equal(t, acc.Owner, randAccount.Owner)
	require.Equal(t, acc.Balance, randAccount.Balance)
	require.Equal(t, acc.Currency, randAccount.Currency)
	require.Equal(t, acc.ID, randAccount.ID)
	require.WithinDuration(t, acc.CreatedAt, randAccount.CreatedAt, time.Second)

}

func TestQueries_UpdateAccount(t *testing.T) {
	randAccount := createFakeAccount(t)
	arg := UpdateAccountParams{ID: randAccount.ID,
		Balance: createRandomInteger(1000, 100000)}
	acc, er := queryTest.UpdateAccount(context.Background(), arg)

	require.NoError(t, er)
	require.NotEmpty(t, acc)
	require.Equal(t, acc.Owner, randAccount.Owner)
	require.Equal(t, acc.Currency, randAccount.Currency)
	require.Equal(t, acc.ID, randAccount.ID)
	require.WithinDuration(t, acc.CreatedAt, randAccount.CreatedAt, time.Second)
	require.Equal(t, acc.Balance, arg.Balance)
}

func TestQueries_DeleteAccount(t *testing.T) {
	randAccount := createFakeAccount(t)

	err := queryTest.DeleteAccount(context.Background(), randAccount.ID)
	require.NoError(t, err)

	acc, err := queryTest.GetAccount(context.Background(), randAccount.ID)

	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, acc)
}

func TestQueries_ListAccounts(t *testing.T) {

	for i := 0; i < 10; i++ {
		createFakeAccount(t)
	}
	arg := ListAccountsParams{Limit: 5, Offset: 5}

	accs, err := queryTest.ListAccounts(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, accs)

	require.Len(t, accs, 5)

	for _, acc := range accs {
		require.NotEmpty(t, acc)
	}
}
