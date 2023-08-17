package db

import (
	"context"
	"simple_bank/util"
	"testing"

	"github.com/stretchr/testify/require"
)

func createRandomEntry(t *testing.T, account Account) Entry {
	args := CreateEntryParams{
		AccountID: account.ID,
		Amount:    util.RandomMoney(),
	}

	entry, err := testQueries.CreateEntry(context.Background(), args)
	require.NoError(t, err)
	require.NotEmpty(t, entry)

	require.Equal(t, args.AccountID, entry.AccountID)
	require.Equal(t, args.Amount, entry.Amount)

	require.NotZero(t, entry.ID)
	require.NotZero(t, entry.CreatedAt)

	paramUpdateAccount := UpdateAccountParams{
		ID:      account.ID,
		Balance: account.Balance + args.Amount,
	}

	accountUpdated, err1 := testQueries.UpdateAccount(context.Background(), paramUpdateAccount)

	require.NoError(t, err1)
	require.NotEmpty(t, accountUpdated)

	require.Equal(t, accountUpdated.Balance, account.Balance+args.Amount)

	return entry
}

func TestCreateEntry(t *testing.T) {
	account := createRandomAccount(t)
	createRandomEntry(t, account)
}

func TestGetEntry(t *testing.T) {
	account := createRandomAccount(t)
	entry := createRandomEntry(t, account)
	entryExists, err := testQueries.GetEntry(context.Background(), entry.ID)

	require.NoError(t, err)
	require.NotEmpty(t, entryExists)

	require.Equal(t, entry.ID, entryExists.ID)
	require.Equal(t, entry.AccountID, entryExists.AccountID)
	require.Equal(t, entry.Amount, entryExists.Amount)
}

func TestListEntries(t *testing.T) {
	account := createRandomAccount(t)
	for i := 0; i < 10; i++ {
		createRandomEntry(t, account)
	}

	arg := ListEntriesParams{
		AccountID: account.ID,
		Limit:     5,
		Offset:    5,
	}

	entries, err := testQueries.ListEntries(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, entries, 5)

	for _, entry := range entries {
		require.NotEmpty(t, entry)
	}
}
