package test

import (
	"testing"

	"github.com/pkg/errors"
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

func TestFunding_GetFundingPayments(t *testing.T) {

	ftx := prepForTest(t)

	payments, err := ftx.Funding.GetFundingPayments(nil,nil,nil)
	if err != nil {
		t.Fatal(errors.WithStack(err))
	}
	if i, p := range payments {
		if i > 9 {
			return
		}
		t.Logf("Payment: %+v\n", *p)
	}
}
