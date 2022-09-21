package testfutures

import (
	"testing"
	"time"

	"github.com/pkg/errors"
	"github.com/sanjujosh/go-ftx/api"
	"github.com/sanjujosh/go-ftx/models"
)

const N = 9

func TestFutures_GetFutures(t *testing.T) {

	ftx := api.New()

	futures, err := ftx.Futures.GetFutures()
	if err != nil {
		t.Fatal(err)
	}

	cnt := 0

	for _, p := range futures {
		if cnt > N {
			return
		}
		if p.Underlying != "BTC" {
			continue
		}
		t.Logf("Future: %+v\n", *p)
		cnt++
	}
}

const fut string = "YFII-PERP"

func TestFutures_GetFutureByName(t *testing.T) {

	ftx := api.New()

	future := models.Future{}
	err := ftx.Futures.GetFutureByName(fut, &future)
	if err != nil {
		t.Fatal(errors.WithStack(err))
	}
	t.Logf("Future: %+v\n", future)
}

func TestFutures_GetFutureStats(t *testing.T) {

	ftx := api.New()
	stats := models.FutureStats{}
	err := ftx.Futures.GetFutureStats(fut, &stats)
	if err != nil {
		t.Fatal(errors.WithStack(err))
	}
	t.Logf("%s stats: %+v\n", fut, stats)
}

func TestFutures_GetFundingRates(t *testing.T) {

	ftx := api.New()

	rates, err := ftx.Futures.GetFundingRates()
	if err != nil {
		t.Fatal(errors.WithStack(err))
	}
	for i, p := range rates {
		if i > 9 {
			return
		}
		t.Logf("Rates: %+v\n", *p)
	}

}

func TestFutures_GetIndexWeights(t *testing.T) {

	/*
		ftx := api.New()
		index := "BTC-PERP"

		weights, err := ftx.Futures.GetIndexWeights(index)
		if err != nil {
			t.Fatal(errors.WithStack(err))
		}
		t.Logf("Weights: %+v\n", weights)
	*/
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

	index := "BTC"
	resolution := 60
	start := time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC).Unix()
	end := start + 60

	histIndex, err := ftx.Futures.GetHistoricalIndex(
		index,
		&models.HistoricalIndexParams{
			Resolution: api.PtrInt(resolution),
			StartTime:  &start,
			EndTime:    &end,
		})

	if err != nil {
		t.Fatal(errors.WithStack(err))
	}

	t.Logf("Length: %d", len(histIndex))

	for _, p := range histIndex {
		if p == nil {
			t.Fatal("nil pointer")
		}
		t.Logf("Historical Index: %+v\n", *p)
	}
}
