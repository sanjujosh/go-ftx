package testacct

import (
	"os"
	"testing"

	"github.com/sanjujosh/go-ftx/api"
	"github.com/sanjujosh/go-ftx/models"
	"github.com/shopspring/decimal"
)

func TestAccount_GetAccountInformation(t *testing.T) {

	ftx := api.New(
		api.WithAuth(os.Getenv("FTX_PROD_MAIN_KEY"), os.Getenv("FTX_PROD_MAIN_SECRET")),
	)
	err := ftx.SetServerTimeDiff()
	if err != nil {
		t.Fatal(err)
	}

	account := models.AccountInformation{}
	err = ftx.Account.GetAccountInformation(&account)
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("Account: %+v\n", account)
}

func TestAccount_GetPositions(t *testing.T) {

	ftx := api.New(
		api.WithAuth(os.Getenv("FTX_PROD_MAIN_KEY"), os.Getenv("FTX_PROD_MAIN_SECRET")),
	)
	err := ftx.SetServerTimeDiff()
	if err != nil {
		t.Fatal(err)
	}

	positions, err := ftx.Account.GetPositions()
	if err != nil {
		t.Fatal(err)
	}

	for _, p := range positions {
		t.Logf("Position: %+v\n", *p)
	}
}

func TestAccount_ChangeAccountLeverage(t *testing.T) {

	ftx := api.New(
		api.WithAuth(os.Getenv("FTX_PROD_MAIN_KEY"), os.Getenv("FTX_PROD_MAIN_SECRET")),
	)
	err := ftx.SetServerTimeDiff()
	if err != nil {
		t.Fatal(err)
	}

	leverage := 10.0

	result, err := ftx.Account.ChangeAccountLeverage(leverage)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("Result: %s\n", result)

	account := models.AccountInformation{}
	err = ftx.Account.GetAccountInformation(&account)

	if err != nil {
		t.Fatal(err)
	}

	l := decimal.NewFromFloat(leverage)
	if !l.Equal(account.Leverage) {
		t.Fatalf("Account leverage not equal to desired leverage: %v, %v", account.Leverage, l)
	}
}
