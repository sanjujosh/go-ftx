package testwallet

import (
	"os"
	"testing"

	"github.com/pkg/errors"
	"github.com/uscott/go-ftx/api"
	"github.com/uscott/go-ftx/models"
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

func TestWallet_GetCoins(t *testing.T) {

	ftx := prepForTest(t)

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

func TestWallet_GetBalances(t *testing.T) {

	ftx := prepForTest(t)

	balances, err := ftx.Wallet.GetBalances()
	if err != nil {
		t.Fatal(errors.WithStack(err))
	}
	for i, bal := range balances {
		if i > 9 {
			return
		}
		t.Logf("Balance: %+v\n", *bal)
	}
}

func TestWallet_GetDepositAddress(t *testing.T) {

	ftx := prepForTest(t)

	address, tag, err := ftx.Wallet.GetDepositAddress("BTC", nil)
	if err != nil {
		t.Fatal(errors.WithStack(err))
	}
	t.Logf("Address, Tag: %s, %s\n", address, tag)
}

func TestWallet_GetDepositHistory(t *testing.T) {

	ftx := prepForTest(t)

	hist, err := ftx.Wallet.GetDepositHistory(&models.DepositHistoryParams{})
	if err != nil {
		t.Fatal(errors.WithStack(err))
	}
	for i, h := range hist {
		if i > 9 {
			return
		}
		t.Logf("Deposit: %+v\n", *h)
	}
}

func TestWallet_GetWithdrawalHistory(t *testing.T) {

	ftx := prepForTest(t)

	hist, err := ftx.Wallet.GetWithdrawalHistory(&models.WithdrawalHistoryParams{})
	if err != nil {
		t.Fatal(errors.WithStack(err))
	}
	for i, h := range hist {
		if i > 9 {
			return
		}
		t.Logf("Withdrawal: %+v\n", *h)
	}
}

func TestWallet_GetAirdrops(t *testing.T) {

	ftx := prepForTest(t)

	drops, err := ftx.Wallet.GetAirdrops(&models.AirDropParams{})
	if err != nil {
		t.Fatal(errors.WithStack(err))
	}
	for i, d := range drops {
		if i > 9 {
			return
		}
		t.Logf("Drop: %+v\n", *d)
	}
}

func TestWallet_GetSavedAddresses(t *testing.T) {

	ftx := prepForTest(t)

	addresses, err := ftx.Wallet.GetSavedAddresses(api.PtrString("BTC"))
	if err != nil {
		t.Fatal(errors.WithStack(err))
	}
	for i, a := range addresses {
		if i > 9 {
			return
		}
		t.Logf("Address: %+v\n", *a)
	}
}
