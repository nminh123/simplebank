package db

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/nminh123/simplebank/util"
	"github.com/stretchr/testify/require"
)

func createRandomEntry(t *testing.T) Entry {
	account := createRandomAccount(t)
	arg := CreateEntryParams{
		AccountID: account.ID,
		Amount:    util.RandomMoney(),
	}

	entr, err := testQueries.CreateEntry(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, entr)

	require.Equal(t, arg.AccountID, entr.AccountID)
	require.Equal(t, arg.Amount, entr.Amount)

	require.NotZero(t, entr.ID)
	require.NotZero(t, entr.CreatedAt)

	return entr
}

func TestCreateEntry(t *testing.T) {
	createRandomEntry(t)
}

func TestGetEntry(t *testing.T) {
	ent1 := createRandomEntry(t)
	ent2, err := testQueries.GetEntry(context.Background(), ent1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, ent2)

	require.Equal(t, ent1.ID, ent2.ID)
	require.Equal(t, ent1.AccountID, ent2.AccountID)
	require.Equal(t, ent1.Amount, ent2.Amount)
	require.Equal(t, ent1.CreatedAt, ent2.CreatedAt)
	require.WithinDuration(t, ent1.CreatedAt, ent2.CreatedAt, time.Second)
}

func TestUpdateEntry(t *testing.T) {
	ent1 := createRandomEntry(t)

	arg := UpdateEntryParams{
		ID:     ent1.ID,
		Amount: ent1.Amount,
	}
	ent2, err := testQueries.UpdateEntry(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, ent2)

	require.Equal(t, ent1.ID, ent2.ID)
	require.Equal(t, ent1.AccountID, ent2.AccountID)
	require.Equal(t, arg.Amount, ent2.Amount)
	require.Equal(t, ent1.CreatedAt, ent2.CreatedAt)
	require.WithinDuration(t, ent1.CreatedAt, ent2.CreatedAt, time.Second)
}

func TestDeleteEntry(t *testing.T) {
	acc1 := createRandomEntry(t)
	err := testQueries.DeleteEntry(context.Background(), acc1.ID)
	require.NoError(t, err)

	acc2, err := testQueries.GetEntry(context.Background(), acc1.ID)
	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, acc2)
}

func TestListEntries(t *testing.T) {
	for i := 0; i < 10; i++ {
		createRandomEntry(t)
	}

	arg := ListEntriesParams{
		Limit:  5,
		Offset: 5,
	}

	entrs, err := testQueries.ListEntries(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, entrs, 5)
	for _, entr := range entrs {
		require.NotEmpty(t, entr)
	}
}
