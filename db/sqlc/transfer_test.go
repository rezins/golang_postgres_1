package db

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func createRandomTransfer(t *testing.T, accountFrom Account, accountTo Account) Transfer {
	args := CreateTransferParams{
		FromAccountID: accountFrom.ID,
		ToAccountID:   accountTo.ID,
		Amount:        50,
	}

	transfer, err := testQueries.CreateTransfer(context.Background(), args)
	require.NoError(t, err)
	require.NotEmpty(t, transfer)

	require.Equal(t, transfer.FromAccountID, accountFrom.ID)
	require.Equal(t, transfer.ToAccountID, accountTo.ID)
	require.Equal(t, transfer.Amount, int64(50))

	require.NotZero(t, transfer.ID)
	require.NotZero(t, transfer.CreatedAt)

	paramUpdateAccountFrom := UpdateAccountParams{
		ID:      accountFrom.ID,
		Balance: accountFrom.Balance - 50,
	}

	accountFromUpdated, err1 := testQueries.UpdateAccount(context.Background(), paramUpdateAccountFrom)

	require.NoError(t, err1)
	require.NotEmpty(t, accountFromUpdated)

	require.Equal(t, accountFromUpdated.Balance, accountFrom.Balance-50)

	paramUpdateAccountTo := UpdateAccountParams{
		ID:      accountTo.ID,
		Balance: accountTo.Balance + 50,
	}

	accountToUpdated, err2 := testQueries.UpdateAccount(context.Background(), paramUpdateAccountTo)

	require.NoError(t, err2)
	require.NotEmpty(t, accountToUpdated)

	require.Equal(t, accountToUpdated.Balance, accountTo.Balance+50)

	return transfer
}

func TestCreateTransfer(t *testing.T) {
	accountFrom := createRandomAccount(t)
	accountTo := createRandomAccount(t)
	createRandomTransfer(t, accountFrom, accountTo)

}

func TestGetTransfer(t *testing.T) {
	accountFrom := createRandomAccount(t)
	accountTo := createRandomAccount(t)
	transfer := createRandomTransfer(t, accountFrom, accountTo)
	transferExists, err := testQueries.GetTransfer(context.Background(), transfer.ID)

	require.NoError(t, err)
	require.NotEmpty(t, transferExists)

	require.Equal(t, transfer.ID, transferExists.ID)
	require.Equal(t, transfer.FromAccountID, transferExists.FromAccountID)
	require.Equal(t, transfer.ToAccountID, transferExists.ToAccountID)
	require.Equal(t, transfer.Amount, transferExists.Amount)
}

func TestListTransfer(t *testing.T) {
	accountFrom := createRandomAccount(t)
	accountTo := createRandomAccount(t)

	for i := 0; i < 5; i++ {
		createRandomTransfer(t, accountFrom, accountTo)
	}

	arg := ListTransfersParams{
		FromAccountID: accountFrom.ID,
		ToAccountID:   accountFrom.ID,
		Limit:         3,
		Offset:        0,
	}

	transfers, err := testQueries.ListTransfers(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, transfers, 3)

	for _, transfer := range transfers {
		require.NotEmpty(t, transfer)
	}

}
