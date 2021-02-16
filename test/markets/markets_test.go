package testmarkets

import (
	"testing"
	"time"

	"github.com/uscott/go-ftx/api"
	"github.com/uscott/go-ftx/models"
)

const N int = 9

var (
	PtrInt   = api.PtrInt
	PtrInt64 = api.PtrInt64
)

func TestMarkets_GetMarkets(t *testing.T) {
	ftx := api.New()

	markets, err := ftx.Markets.GetMarkets()
	if err != nil {
		t.Fatal(err)
	}

	cnt := 0
	for _, p := range markets {
		if cnt > N {
			return
		}
		cnt++
		t.Logf("%s market: %+v\n", p.Name, *p)
	}
}

func TestMarkets_GetMarketByName(t *testing.T) {

	ftx := api.New()
	market := models.Market{}

	expected := &models.Market{
		Name:          "BTC/USD",
		BaseCurrency:  "BTC",
		QuoteCurrency: "USD",
		Enabled:       true,
	}

	if err := ftx.Markets.GetMarketByName(expected.Name, &market); err != nil {
		t.Fatal(err)
	}
	t.Logf("Market: %+v", market)

	if ftx.Markets.GetMarketByName("incorrect", &market) == nil {
		t.Fatal("Should have gotten an error")
	}

	if err := ftx.Markets.GetMarketByName("BTC-PERP", &market); err != nil {
		t.Fatal(err)
	}
	t.Logf("\n%s market: %+v\n", market.Name, market)

}

func TestMarkets_GetOrderBook(t *testing.T) {

	ftx := api.New()
	ob := models.OrderBook{}

	if err := ftx.Markets.GetOrderBook("ETH/BTC", nil, &ob); err != nil {
		t.Fatal(err)
	}

	depth := 30
	if err := ftx.Markets.GetOrderBook("ETH/BTC", &depth, &ob); err != nil {
		t.Fatal(err)
	}
	if len(ob.Asks) != depth || len(ob.Bids) != depth {
		t.Fatalf("Lengths are wrong: %d, %d, %d", depth, len(ob.Asks), len(ob.Bids))
	}

	if ftx.Markets.GetOrderBook("failed", &depth, &ob) == nil {
		t.Fatal("Should have gotten an err")
	}
}

func TestMarkets_GetTrades(t *testing.T) {

	ftx := api.New()

	trades, err := ftx.Markets.GetTrades("ETH/BTC", nil)
	if err != nil {
		t.Fatal(err)
	}
	if trades == nil {
		t.Fatal("Trades should not be nil")
	}

	limit := 10
	trades, err = ftx.Markets.GetTrades("ETH/BTC", &models.GetTradesParams{
		Limit: &limit,
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(trades) != limit {
		t.Fatalf("Trades have wrong length: %d, %d", limit, len(trades))
	}

	trades, err = ftx.Markets.GetTrades("ETH/BTC", &models.GetTradesParams{
		Limit:     &limit,
		StartTime: PtrInt64(time.Now().Add(-5 * time.Hour).Unix()),
		EndTime:   PtrInt64(time.Now().Unix()),
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(trades) != limit {
		t.Fatalf("Trades have wrong length: %d, %d", limit, len(trades))
	}
}

func TestMarkets_GetHistoricalPrices(t *testing.T) {

	ftx := api.New()

	prices, err := ftx.Markets.GetHistoricalPrices("ETH/BTC", nil)
	if err == nil {
		t.Fatal("Should have gotten an error")
	}
	if prices != nil {
		t.Fatal("Array should be nil")
	}

	prices, err = ftx.Markets.GetHistoricalPrices(
		"ETH/BTC", &models.GetHistoricalPricesParams{
			Resolution: models.Minute,
		})
	if err != nil {
		t.Fatal(err)
	}
	if prices == nil {
		t.Fatal("Prices should not be nil")
	}

	prices, err = ftx.Markets.GetHistoricalPrices(
		"ETH/BTC", &models.GetHistoricalPricesParams{
			Resolution: models.Minute,
			Limit:      PtrInt(10),
			StartTime:  PtrInt(int(time.Now().Add(-5 * time.Hour).Unix())),
			EndTime:    PtrInt(int(time.Now().Unix())),
		})
	if err != nil {
		t.Fatal(err)
	}
	if prices == nil {
		t.Fatal("Prices should not be nil")
	}
	if len(prices) != 10 {
		t.Fatalf("Length of prices should be 10: %d", len(prices))
	}
}
