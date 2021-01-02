package test

import (
	"fmt"
	"testing"
	"time"

	"github.com/pkg/errors"
	"github.com/uscott/go-ftx/api"
	"github.com/uscott/go-ftx/models"
)

func TestFutures_GetFutures(t *testing.T) {

	ftx := api.New()

	futures, err := ftx.Futures.GetFutures()
	if err != nil {
		t.Fatal(err)
	}
	for _, p := range futures {
		if p == nil {
			t.Fatal("nil pointer")
		}
		fmt.Printf("Description: %s\n", p.Description)
		fmt.Printf("Expiration:  %+v\n", p.Expiry.Format(time.RFC3339))
		fmt.Printf("Name:        %s\n", p.Name)
	}
}

const futName string = "Bitcoin March 2019 Futures"

func TestFutures_GetFutureByName(t *testing.T) {

	ftx := api.New()

	future, err := ftx.Futures.GetFutureByName(futName)
	if err != nil {
		t.Fatal(errors.WithStack(err))
	}
	t.Logf("Future: %+v\n", future)
}

func TestFutures_GetFutureStats(t *testing.T) {

	ftx := api.New()
	stats, err := ftx.Futures.GetFutureStats(futName)
	if err != nil {
		t.Fatal(errors.WithStack(err))
	}
	t.Logf("Stats: %+v\n", stats)
}

func TestFutures_GetFundingRates(t *testing.T) {

	ftx := api.New()
	now := time.Now()
	rates, err := ftx.Futures.GetFundingRates(&models.FundingRatesParams{
		StartTime: api.PtrInt(int(now.Add(-5 * time.Hour).Unix())),
		EndTime:   api.PtrInt(int(now.Unix())),
		Future:    api.PtrString("BTC-PERP"),
	})
	if err != nil {
		t.Fatal(errors.WithStack(err))
	}
	t.Logf("Rates: %+v\n", rates)
}
