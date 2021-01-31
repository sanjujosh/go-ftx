package testspotmargin

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

func TestSpotMargin_GetBorrowRates(t *testing.T) {

	ftx := prepForTest(t)

	rates, err := ftx.SpotMargin.GetBorrowRates()
	if err != nil {
		t.Fatal(errors.WithStack(err))
	}
	for i, r := range rates {
		if i > 9 {
			return
		}
		t.Logf("Rate: %+v\n", *r)
	}
}

func TestSpotMargin_GetLendingRates(t *testing.T) {

	ftx := prepForTest(t)

	rates, err := ftx.SpotMargin.GetLendingRates()
	if err != nil {
		t.Fatal(errors.WithStack(err))
	}
	for i, r := range rates {
		if i > 9 {
			return
		}
		t.Logf("Rate: %+v\n", *r)
	}
}

func TestSpotMargin_GetBorrowSummary(t *testing.T) {

	ftx := prepForTest(t)

	summary, err := ftx.SpotMargin.GetBorrowSummary()
	if err != nil {
		t.Fatal(errors.WithStack(err))
	}
	for i, s := range summary {
		if i > 9 {
			return
		}
		t.Logf("Summary: %+v\n", *s)
	}
}

func TestSpotMargin_GetMarketInfo(t *testing.T) {

	ftx := prepForTest(t)

	info, err := ftx.SpotMargin.GetMarketInfo("BTC/USD")
	if err != nil {
		t.Fatal(errors.WithStack(err))
	}
	t.Logf("Info: %+v\n", *info)
}

func TestSpotMargin_GetBorrowHistory(t *testing.T) {

	ftx := prepForTest(t)

	hist, err := ftx.SpotMargin.GetBorrowHistory()
	if err != nil {
		t.Fatal(errors.WithStack(err))
	}
	for i, h := range hist {
		if i > 9 {
			return
		}
		t.Logf("Hist: %+v\n", *h)
	}
}

func TestSpotMargin_GetLendingHistory(t *testing.T) {

	ftx := prepForTest(t)

	hist, err := ftx.SpotMargin.GetLendingHistory()
	if err != nil {
		t.Fatal(errors.WithStack(err))
	}
	for i, h := range hist {
		if i > 9 {
			return
		}
		t.Logf("Hist: %+v\n", *h)
	}
}

func TestSpotMargin_GetLendingOffers(t *testing.T) {

	ftx := prepForTest(t)

	offers, err := ftx.SpotMargin.GetLendingOffers()
	if err != nil {
		t.Fatal(errors.WithStack(err))
	}
	for i, o := range offers {
		if i > 9 {
			return
		}
		t.Logf("Lending offer: %+v\n", *o)
	}
}

func TestSpotMargin_GetLendingInfo(t *testing.T) {

	ftx := prepForTest(t)

	info, err := ftx.SpotMargin.GetLendingInfo()
	if err != nil {
		t.Fatal(errors.WithStack(err))
	}
	for i, x := range info {
		if i > 9 {
			return
		}
		t.Logf("Info: %+v\n", *x)
	}
}
