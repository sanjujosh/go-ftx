package test

import (
	"os"
	"testing"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/uscott/go-ftx"
)

func TestSubAccounts_CRUD(t *testing.T) {
	godotenv.Load()

	ftx := ftx.New(
		ftx.WithAuth(os.Getenv("FTX_PROD_MAIN_KEY"), os.Getenv("FTX_PROD_MAIN_SECRET")),
	)
	err := ftx.SetServerTimeDiff()
	require.NoError(t, err)

	nickname := "testSubAccount"
	newNickname := "newTestSubAccount"

	t.Run("getAll", func(t *testing.T) {
		subs, err := ftx.SubAccounts.GetSubaccounts()
		assert.NoError(t, err)
		assert.NotNil(t, subs)
		for _, sub := range subs {
			require.True(t, sub.Nickname != nickname)
			require.True(t, sub.Nickname != newNickname)
		}
	})

	t.Run("create", func(t *testing.T) {
		sub, err := ftx.SubAccounts.CreateSubaccount(nickname)
		assert.NoError(t, err)
		assert.NotNil(t, sub)
		assert.Equal(t, nickname, sub.Nickname)
		assert.True(t, sub.Deletable)
		assert.True(t, sub.Editable)
	})

	t.Run("get_balances", func(t *testing.T) {
		balances, err := ftx.SubAccounts.GetSubaccountBalances(nickname)
		assert.NoError(t, err)
		assert.NotNil(t, balances)
	})

	t.Run("update", func(t *testing.T) {
		err = ftx.SubAccounts.ChangeSubaccount(nickname, newNickname)
		assert.NoError(t, err)
	})

	t.Run("delete", func(t *testing.T) {
		err = ftx.SubAccounts.DeleteSubaccount(newNickname)
		assert.NoError(t, err)
	})

	t.Run("check", func(t *testing.T) {
		subs, err := ftx.SubAccounts.GetSubaccounts()
		assert.NoError(t, err)
		assert.NotNil(t, subs)
		for _, sub := range subs {
			require.True(t, sub.Nickname != nickname)
			require.True(t, sub.Nickname != newNickname)
		}
	})
}
