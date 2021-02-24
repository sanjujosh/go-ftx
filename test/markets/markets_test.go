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
		Name:          "YFII/USD",
		BaseCurrency:  "YFII",
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
	symbol := "BTC/USD"

	trades, err := ftx.Markets.GetTrades(symbol, nil)
	if err != nil {
		t.Fatal(err)
	}
	limit := 10
	trades, err = ftx.Markets.GetTrades(symbol, &models.GetTradesParams{
		Limit: &limit,
	})
	if err != nil {
		t.Fatal(err)
	}

	trades, err = ftx.Markets.GetTrades(symbol, &models.GetTradesParams{
		Limit:     &limit,
		StartTime: PtrInt64(time.Now().Add(-5 * time.Hour).Unix()),
		EndTime:   PtrInt64(time.Now().Unix()),
	})
	if err != nil {
		t.Fatal(err)
	}
	for i, x := range trades {
		if i > 5 {
			break
		}
		t.Logf("Trade: %+v", *x)
	}
}

func TestMarkets_GetHistoricalPrices(t *testing.T) {

	ftx := api.New()
	symbol := "LEO/USD"
	limit := 10

	params := &models.GetHistoricalPricesParams{
		Resolution: 24 * models.Hour,
		Limit: &limit,
	}

	prices, err := ftx.Markets.GetHistoricalPrices(symbol, params)
	if err != nil {
		t.Fatal(err)
	}
	if prices == nil {
		t.Fatal("Prices should not be nil")
	}
	for i, p := range prices {
		t.Logf("Historical price: %+v", *p)
		if i > 5 {
			break
		}
	}

	now := time.Now()

	params.EndTime = PtrInt(int(now.Unix()))
	params.StartTime = PtrInt(int(now.Add(-7 * 24 * time.Hour).Unix()))

	prices, err = ftx.Markets.GetHistoricalPrices(symbol, params)
	if err != nil {
		t.Fatal(err)
	}
	if prices == nil {
		t.Fatal("Prices should not be nil")
	}
	for i, p := range prices {
		t.Logf("Historical price: %+v", *p)
		if i > 10 {
			break
		}
	}
}
