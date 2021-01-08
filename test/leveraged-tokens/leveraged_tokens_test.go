package test

import (
	"os"
	"testing"

	"github.com/pkg/errors"
	"github.com/uscott/go-ftx/api"
)

func prepForTest(t *testing.T) *api.Client {

	ftx := api.New(
		api.WithAuth(os.Getenv("FTX_PROD_MAIN_KEY"), os.Getenv("FTX_PROD_MAIN_SECRET")),
	)
	if err := ftx.SetServerTimeDiff(); err != nil {
		t.Fatal(errors.WithStack(err))
	}
	return ftx
}

func TestLeveragedTokens_ListLeveragedTokens(t *testing.T) {

	ftx := prepForTest(t)

	list, err := ftx.LeveragedTokens.ListLeveragedTokens()
	if err != nil {
		t.Fatal(errors.WithStack(err))
	}
	for i, l := range list {
		if i > 9 {
			return
		}
		t.Logf("Token: %+v\n", *l)
	}
}

func TestLeveragedTokens_GetTokenInfo(t *testing.T) {

}

func TestLeveragedTokens_GetLeveragedTokenBalances(t *testing.T) {

	ftx := prepForTest(t)

	balances, err := ftx.LeveragedTokens.GetLeveragedTokenBalances()
	if err != nil {
		t.Fatal(errors.WithStack(err))
	}
	for i, b := range balances {
		if i > 9 {
			return
		}
		t.Logf("Balance: %+v\n", *b)
	}
}

func TestLeveragedTokens_ListLeveragedTokenCreationRequests(t *testing.T) {

	ftx := prepForTest(t)

	requests, err := ftx.LeveragedTokens.ListLeveragedTokenCreationRequests()
	if err != nil {
		t.Fatal(errors.WithStack(err))
	}
	for i, req := range requests {
		if i > 9 {
			return
		}
		t.Logf("Request: %+v\n", *req)
	}
}

func TestLeveragedTokens_RequestLeveragedTokenCreation(t *testing.T) {

}

func TestLeveragedTokens_ListLeveragedTokenRedemptionRequests(t *testing.T) {

	ftx := prepForTest(t)

	list, err := ftx.LeveragedTokens.ListLeveragedTokenRedemptionRequests()
	if err != nil {
		t.Fatal(errors.WithStack(err))
	}
	for i, l := range list {
		if i > 9 {
			return
		}
		t.Logf("Request: %+v\n", *l)
	}
}
