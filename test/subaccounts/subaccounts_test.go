package test

import (
	"os"
	"testing"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/uscott/go-ftx/api"
	"github.com/uscott/go-ftx/models"
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

	subs, err := ftx.SubAccounts.GetSubaccounts()
	assert.NoError(t, err)
	assert.NotNil(t, subs)
	for _, sub := range subs {
		t.Logf("Subaccount: %+v\n", *sub)
		if sub.Nickname == nickname || sub.Nickname == newNickname {
			_, err = ftx.SubAccounts.DeleteSubaccount(sub.Nickname)
			if err != nil {
				t.Fatal(err)

			}
		}
	}

	sub, err := ftx.SubAccounts.CreateSubaccount(nickname)
	assert.NoError(t, err)
	assert.NotNil(t, sub)
	assert.Equal(t, nickname, sub.Nickname)
	assert.True(t, sub.Deletable)
	assert.True(t, sub.Editable)
	t.Logf("Subaccount: %+v\n", *sub)

	ftx.SubAccount = api.PtrString(nickname)
	account := &models.AccountInformation{}
	if err = ftx.Account.GetAccountInformation(account); err != nil {
		t.Fatal(err)
	}
	t.Logf("%s account: %+v", nickname, *account)

	positions, err := ftx.Account.GetPositions()
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%s positions:\n", nickname)
	for _, p := range positions {
		t.Logf("%+v\n", *p)
	}

	ftx.SubAccount = nil

	balances, err := ftx.SubAccounts.GetSubaccountBalances(nickname)
	assert.NoError(t, err)
	assert.NotNil(t, balances)
	for _, bal := range balances {
		t.Logf("Balance: %+v\n", *bal)
	}

	result, err := ftx.SubAccounts.ChangeSubaccount(nickname, newNickname)
	assert.NoError(t, err)
	t.Logf("Update result: %+v\n", result)

	result, err = ftx.SubAccounts.DeleteSubaccount(newNickname)
	assert.NoError(t, err)
	t.Logf("Delete result: %+v\n", result)

	subs, err = ftx.SubAccounts.GetSubaccounts()
	assert.NoError(t, err)
	assert.NotNil(t, subs)
	for _, sub := range subs {
		require.True(t, sub.Nickname != nickname)
		require.True(t, sub.Nickname != newNickname)
		t.Logf("Check subaccount: %+v\n", *sub)
	}
}
