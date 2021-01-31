package test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

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
	assert.NoError(t, err)
	assert.NotNil(t, markets)

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
	req := require.New(t)
	market := models.Market{}

	t.Run("success", func(t *testing.T) {
		expected := &models.Market{
			Name:          "BTC/USD",
			BaseCurrency:  "BTC",
			QuoteCurrency: "USD",
			Enabled:       true,
		}

		err := ftx.Markets.GetMarketByName(expected.Name, &market)
		req.NoError(err)
		req.NotNil(&market)
		req.Equal(expected.Name, market.Name)
		req.Equal(expected.BaseCurrency, market.BaseCurrency)
		req.Equal(expected.QuoteCurrency, market.QuoteCurrency)
		req.Equal(expected.Enabled, market.Enabled)
		t.Logf("\n%s market: %+v\n", market.Name, market)
	})

	t.Run("not_found", func(t *testing.T) {
		expected := &models.Market{
			Name:          "incorrect",
			BaseCurrency:  "ETH",
			QuoteCurrency: "BTC",
			Enabled:       true,
		}

		err := ftx.Markets.GetMarketByName(expected.Name, &market)
		req.Error(err)
	})

	if err := ftx.Markets.GetMarketByName("BTC-PERP", &market); err != nil {
		t.Fatal(err)
	}
	t.Logf("\n%s market: %+v\n", market.Name, market)

}

func TestMarkets_GetOrderBook(t *testing.T) {

	ftx := api.New()

	req := require.New(t)
	ob := models.OrderBook{}

	t.Run("success", func(t *testing.T) {
		err := ftx.Markets.GetOrderBook("ETH/BTC", nil, &ob)
		req.NoError(err)
		req.NotNil(&ob)
	})

	t.Run("success_with_depth", func(t *testing.T) {
		depth := 30
		err := ftx.Markets.GetOrderBook("ETH/BTC", &depth, &ob)
		req.NoError(err)
		req.NotNil(&ob)
		req.Len(ob.Asks, depth)
		req.Len(ob.Bids, depth)
	})

	t.Run("failed_market", func(t *testing.T) {
		depth := 30
		err := ftx.Markets.GetOrderBook("failed", &depth, &ob)
		req.Error(err)
	})
}

func TestMarkets_GetTrades(t *testing.T) {

	ftx := api.New()
	req := require.New(t)

	t.Run("success", func(t *testing.T) {
		trades, err := ftx.Markets.GetTrades("ETH/BTC", nil)
		req.NoError(err)
		req.NotNil(trades)
	})

	t.Run("success_with_limit", func(t *testing.T) {
		limit := 10
		trades, err := ftx.Markets.GetTrades("ETH/BTC", &models.GetTradesParams{
			Limit: &limit,
		})
		req.NoError(err)
		req.NotNil(trades)
		req.Len(trades, limit)
	})

	t.Run("success_with_params", func(t *testing.T) {
		limit := 10
		trades, err := ftx.Markets.GetTrades("ETH/BTC", &models.GetTradesParams{
			Limit:     &limit,
			StartTime: PtrInt64(time.Now().Add(-5 * time.Hour).Unix()),
			EndTime:   PtrInt64(time.Now().Unix()),
		})
		req.NoError(err)
		req.NotNil(trades)
		req.Len(trades, limit)
	})
}

func TestMarkets_GetHistoricalPrices(t *testing.T) {

	ftx := api.New()
	req := require.New(t)

	t.Run("failed", func(t *testing.T) {
		prices, err := ftx.Markets.GetHistoricalPrices("ETH/BTC", nil)
		req.Error(err)
		req.Nil(prices)
	})

	t.Run("success_with_resolution", func(t *testing.T) {
		prices, err := ftx.Markets.GetHistoricalPrices(
			"ETH/BTC", &models.GetHistoricalPricesParams{
				Resolution: models.Minute,
			})
		req.NoError(err)
		req.NotNil(prices)
	})

	t.Run("success_with_params", func(t *testing.T) {
		prices, err := ftx.Markets.GetHistoricalPrices(
			"ETH/BTC", &models.GetHistoricalPricesParams{
				Resolution: models.Minute,
				Limit:      PtrInt(10),
				StartTime:  PtrInt(int(time.Now().Add(-5 * time.Hour).Unix())),
				EndTime:    PtrInt(int(time.Now().Unix())),
			})
		req.NoError(err)
		req.NotNil(prices)
		req.Len(prices, 10)
	})
}
