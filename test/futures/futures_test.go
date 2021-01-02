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
	for i, p := range futures {
		if p == nil {
			t.Fatal("nil pointer")
		}
		fmt.Printf("Description: %s\n", p.Description)
		fmt.Printf("Expiration:  %+v\n", p.Expiry.Format(time.RFC3339))
		fmt.Printf("Name:        %s\n", p.Name)
		if i > 10 {
			break
		}
	}
}

const futName string = "BTC-PERP"

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
		Future:    api.PtrString(futName),
	})
	if err != nil {
		t.Fatal(errors.WithStack(err))
	}
	for i, p := range rates {
		if p == nil {
			t.Fatal("nil pointer")
		}
		t.Logf("Rates: %+v\n", *p)
		if i > 10 {
			break
		}
	}

}

const index string = "BTC"

func TestFutures_GetIndexWeights(t *testing.T) {

	ftx := api.New()

	weights, err := ftx.Futures.GetIndexWeights(index)
	if err != nil {
		t.Fatal(errors.WithStack(err))
	}
	t.Logf("Weights: %+v\n", weights)
}

func TestFutures_GetExpiredFutures(t *testing.T) {

	ftx := api.New()
	futures, err := ftx.Futures.GetExpiredFutures()
	if err != nil {
		t.Fatal(errors.WithStack(err))
	}
	for i, p := range futures {
		if p == nil {
			t.Fatal("nil pointer")
		}
		t.Logf("Expired Future: %+v\n", *p)
		if i > 10 {
			break
		}
	}
}

func TestFutures_GetHistoricalIndex(t *testing.T) {

	ftx := api.New()
	limit, resolution, now := 30, 60, time.Now()
	histIndex, err := ftx.Futures.GetHistoricalIndex(
		index,
		&models.HistoricalIndexParams{
		Resolution: api.PtrInt(resolution),
		Limit:      api.PtrInt(limit),
		StartTime:  api.PtrInt(int(now.Add(-5 * time.Hour).Unix())),
		EndTime:    api.PtrInt(int(now.Unix())),
	})
	if err != nil {
		t.Fatal(errors.WithStack(err))
	}
	for i, p := range histIndex {
		if p == nil {
			t.Fatal("nil pointer")
		}
		t.Logf("Historical Index: %+v\n", *p)
		if i > 10 {
			break
		}
	}
}
