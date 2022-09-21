package testsubaccounts

import (
	"os"
	"testing"

	"github.com/sanjujosh/go-ftx/api"
	"github.com/sanjujosh/go-ftx/models"
)

func TestSubAccounts_CRUD(t *testing.T) {

	ftx := api.New(
		api.WithAuth(os.Getenv("FTX_PROD_MAIN_KEY"), os.Getenv("FTX_PROD_MAIN_SECRET")),
	)
	err := ftx.SetServerTimeDiff()
	if err != nil {
		t.Fatal(err)
	}

	nickname := "testSubAccount"
	newNickname := "newTestSubAccount"

	subs, err := ftx.SubAccounts.GetSubaccounts()
	if err != nil {
		t.Fatal(err)
	}
	if subs == nil {
		t.Fatal("subs should not be nil")
	}

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
	if err != nil {
		t.Fatal(err)
	}
	if sub == nil {
		t.Fatal("sub should not be nil")
	}
	if sub.Nickname != nickname {
		t.Fatalf("Wrong nickname: %s, %s", sub.Nickname, nickname)
	}
	if !sub.Deletable {
		t.Fatal("Not deletable")
	}
	if !sub.Editable {
		t.Fatal("Not editable")
	}
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
	if err != nil {
		t.Fatal(err)
	}
	if balances == nil {
		t.Fatal("balance should not be nil")
	}

	for _, bal := range balances {
		t.Logf("Balance: %+v\n", *bal)
	}

	result, err := ftx.SubAccounts.ChangeSubaccount(nickname, newNickname)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("Update result: %+v\n", result)

	result, err = ftx.SubAccounts.DeleteSubaccount(newNickname)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("Delete result: %+v\n", result)

	subs, err = ftx.SubAccounts.GetSubaccounts()
	if err != nil {
		t.Fatal(err)
	}

	if subs == nil {
		t.Fatal("subs should not be nil")
	}

	for _, sub := range subs {
		if sub.Nickname == nickname {
			t.Fatal("Wrong nickname: %s, %s", sub.Nickname, newNickname)
		}
		if sub.Nickname != newNickname {
			t.Fatal("Wrong nickname: %s, %s", sub.Nickname, newNickname)
		}
		t.Logf("Check subaccount: %+v\n", *sub)
	}
}
