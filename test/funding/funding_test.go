package testfunding

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"
	"time"

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

func TestFunding_GetFundingPayments(t *testing.T) {

	ftx := prepForTest(t)
	ftx.SubAccount = api.PtrString("SWAPS")

	start := time.Date(2021, 2, 4, 0, 0, 0, 0, time.UTC).Unix()
	end := time.Date(2021, 3, 5, 0, 0, 0, 0, time.UTC).Unix()

	pars := &models.FundingPaymentParams{StartTime: &start, EndTime: &end}
	newpars1, err := api.PrepareQueryParams(pars)
	if err != nil {
		t.Fatal(err)
	}
	b, err := json.Marshal(pars)
	if err != nil {
		t.Fatal(err)
	}
	var newpars2 map[string]interface{}
	if err = json.Unmarshal(b, &newpars2); err != nil {
		t.Fatal(err)
	}
	var newpars3 map[string]string
	newpars3 = make(map[string]string)
	for k, v := range newpars2 {
		newpars3[k] = fmt.Sprintf("%v", v)
	}

	t.Logf("Start = %d", start)
	t.Logf("End   = %d", end)
	t.Logf("pars1 = %+v", newpars1)
	t.Logf("bytes = %+v", string(b))
	t.Logf("pars2 = %+v", newpars2)
	t.Logf("pars3 = %+v", newpars3)

	payments, err := ftx.Funding.GetFundingPayments(nil, &start, &end)
	if err != nil {
		t.Fatal(errors.WithStack(err))
	}

	t.Logf("Length = %d", len(payments))

	for i, p := range payments {
		if i > 1*39 {
			break
		}
		t.Logf("Payment: %+v\n", *p)
	}

	future := api.PtrString("BTMX-PERP")
	payments, err = ftx.Funding.GetFundingPayments(future, &start, &end)
	if err != nil {
		t.Fatal(errors.WithStack(err))
	}
	t.Logf("Length = %d", len(payments))
	for i, p := range payments {
		if i > 1*39 {
			break
		}
		t.Logf("Payment: %+v\n", *p)
	}
}
