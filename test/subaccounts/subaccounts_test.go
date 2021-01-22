package test

import (
	"os"
	"testing"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/uscott/go-ftx/api"
)

func TestSubAccounts_CRUD(t *testing.T) {
	godotenv.Load()

	ftx := api.New(
		api.WithAuth(os.Getenv("FTX_PROD_MAIN_KEY"), os.Getenv("FTX_PROD_MAIN_SECRET")),
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
			t.Logf("Subaccount: %+v\n", *sub)
		}
	})

	t.Run("create", func(t *testing.T) {
		sub, err := ftx.SubAccounts.CreateSubaccount(nickname)
		assert.NoError(t, err)
		assert.NotNil(t, sub)
		assert.Equal(t, nickname, sub.Nickname)
		assert.True(t, sub.Deletable)
		assert.True(t, sub.Editable)
		t.Logf("Subaccount: %+v\n", *sub)
	})

	t.Run("get_balances", func(t *testing.T) {
		balances, err := ftx.SubAccounts.GetSubaccountBalances(nickname)
		assert.NoError(t, err)
		assert.NotNil(t, balances)
		for _, bal := range balances {
			t.Logf("Balance: %+v\n", *bal)
		}
	})

	t.Run("update", func(t *testing.T) {
		result, err := ftx.SubAccounts.ChangeSubaccount(nickname, newNickname)
		assert.NoError(t, err)
		t.Logf("Update result: %+v\n", result)
	})

	t.Run("delete", func(t *testing.T) {
		result, err := ftx.SubAccounts.DeleteSubaccount(newNickname)
		assert.NoError(t, err)
		t.Logf("Delete result: %+v\n", result)
	})

	t.Run("check", func(t *testing.T) {
		subs, err := ftx.SubAccounts.GetSubaccounts()
		assert.NoError(t, err)
		assert.NotNil(t, subs)
		for _, sub := range subs {
			require.True(t, sub.Nickname != nickname)
			require.True(t, sub.Nickname != newNickname)
			t.Logf("Check subaccount: %+v\n", *sub)
		}
	})
}
