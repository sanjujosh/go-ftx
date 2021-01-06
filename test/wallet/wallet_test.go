package test

import (
	"os"
	"testing"

	"github.com/pkg/errors"
	"github.com/uscott/go-ftx/api"
)

func prepForTest() (*api.Client, error) {

	ftx := api.New(
		api.WithAuth(os.Getenv("FTX_PROD_MAIN_KEY"), os.Getenv("FTX_PROD_MAIN_SECRET")),
	)
	if err := ftx.SetServerTimeDiff(); err != nil {
		return nil, errors.WithStack(err)
	}
	return ftx, nil
}

func TestWallet_GetCoins(t *testing.T) {

	ftx, err := prepForTest()
	if err != nil {
		t.Fatal(err)
	}
	coins, err := ftx.Wallet.GetCoins()
	if err != nil {
		t.Fatal(errors.WithStack(err))
	}
	for i, c := range coins {
		if i > 9 {
			return
		}
		t.Logf("Coin: %+v\n", *c)
	}
}
