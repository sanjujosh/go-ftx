package testoptions

import (
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

func TestOptions_ListQuoteRequests(t *testing.T) {

	ftx := prepForTest(t)

	requests, err := ftx.Options.ListQuoteRequests()
	if err != nil {
		t.Fatal(errors.WithStack(err))
	}
	for i, req := range requests {
		if i > 9 {
			return
		}
		t.Logf("Request : %+v\n", *req)
	}
}

func TestOptions_ListUserQuoteRequests(t *testing.T) {

	ftx := prepForTest(t)

	requests, err := ftx.Options.ListUserQuoteRequests()
	if err != nil {
		t.Fatal(errors.WithStack(err))
	}
	for i, req := range requests {
		if i > 9 {
			return
		}
		t.Logf("Request : %+v\n", *req)
	}
}

func TestOptions_CreateQuoteRequest(t *testing.T) {
}

func TestOptions_CancelQuoteRequest(t *testing.T) {}

func TestOptions_GetQuotesForUserQuoteRequest(t *testing.T) {}

func TestOptions_CreateQuote(t *testing.T) {}

func TestOptions_GetUserQuotes(t *testing.T) {

	ftx := prepForTest(t)

	quotes, err := ftx.Options.GetUserQuotes()
	if err != nil {
		t.Fatal(errors.WithStack(err))
	}
	for i, q := range quotes {
		if i > 9 {
			return
		}
		t.Logf("Request : %+v\n", *q)
	}
}

func TestOptions_CancelQuote(t *testing.T) {}

func TestOptions_AcceptQuote(t *testing.T) {}

func TestOptions_GetAccountOptionsInfo(t *testing.T) {

	ftx := prepForTest(t)

	info, err := ftx.Options.GetAccountOptionsInfo()
	if err != nil {
		t.Fatal(errors.WithStack(err))
	}
	t.Logf("Account Options Info: %+v\n", *info)
}

func TestOptions_GetOptionsPositions(t *testing.T) {

	ftx := prepForTest(t)

	positions, err := ftx.Options.GetOptionsPositions()
	if err != nil {
		t.Fatal(errors.WithStack(err))
	}
	for i, p := range positions {
		if i > 9 {
			return
		}
		t.Logf("Position: %+v\n", *p)
	}
}

func TestOptions_GetPublicOptionsTrades(t *testing.T) {

	ftx := prepForTest(t)
	now := time.Now().UTC()
	params := &models.NumberTimeLimit{
		Limit:     api.PtrInt(10),
		StartTime: api.PtrInt64(now.Add(-24 * time.Hour).Unix()),
		EndTime:   api.PtrInt64(now.Unix()),
	}
	trades, err := ftx.Options.GetPublicOptionsTrades(params)
	if err != nil {
		t.Fatal(errors.WithStack(err))
	}
	for i, x := range trades {
		if i > 9 {
			return
		}
		t.Logf("Trade: %+v\n", *x)
	}
}

func TestOptions_GetOptionFills(t *testing.T) {
}

func TestOptions_Get24hOptionsVolume(t *testing.T) {

	ftx := prepForTest(t)
	volume, err := ftx.Options.Get24hOptionVolume()
	if err != nil {
		t.Fatal(errors.WithStack(err))
	}
	t.Logf("24h Volume: %+v\n", *volume)
}

func TestOptions_GetOptionsHistoricalVolumes(t *testing.T) {

	ftx := prepForTest(t)
	now := time.Now().UTC()
	params := &models.NumberTimeLimit{
		Limit:     api.PtrInt(10),
		StartTime: api.PtrInt64(now.Add(-24 * time.Hour).Unix()),
		EndTime:   api.PtrInt64(now.Unix()),
	}
	volumes, err := ftx.Options.GetOptionsHistoricalVolumes(params)
	if err != nil {
		t.Fatal(errors.WithStack(err))
	}
	for i, v := range volumes {
		if i > 9 {
			return
		}
		t.Logf("Volume: %+v\n", *v)
	}
}

func TestOptions_GetOptionsOpenInterest(t *testing.T) {

	ftx := prepForTest(t)

	oi, err := ftx.Options.GetOptionsOpenInterest()
	if err != nil {
		t.Fatal(errors.WithStack(err))
	}
	t.Logf("Open Interest: %+v\n", oi)
}

func TestOptions_GetHistoricalOpenInterest(t *testing.T) {

	ftx := prepForTest(t)
	now := time.Now().UTC()
	params := &models.NumberTimeLimit{
		Limit:     api.PtrInt(10),
		StartTime: api.PtrInt64(now.Add(-24 * time.Hour).Unix()),
		EndTime:   api.PtrInt64(now.Unix()),
	}
	openInterest, err := ftx.Options.GetHistoricalOpenInterest(params)
	if err != nil {
		t.Fatal(errors.WithStack(err))
	}
	for i, oi := range openInterest {
		if i > 9 {
			return
		}
		t.Logf("Open Interest: %+v\n", *oi)
	}
}
