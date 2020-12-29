package test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	ftx "github.com/uscott/go-ftx/api"
	"github.com/uscott/go-ftx/models"
)

const sleepDuration time.Duration = 5 * time.Second

func prepForTest() (*ftx.Client, context.Context, chan struct{}) {
	ftx := ftx.New()
	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan struct{})
	go func() {
		time.Sleep(sleepDuration)
		cancel()
		done <- struct{}{}
	}()
	return ftx, ctx, done
}

func TestStream_SubscribeToTickers(t *testing.T) {

	ftx, ctx, done := prepForTest()

	symbol := "ETH/BTC"
	data, err := ftx.Stream.SubscribeToTickers(ctx, symbol)
	require.NoError(t, err)

	count := 0
	for {
		select {
		case <-done:
			require.True(t, count > 0)
			return
		case msg := <-data:
			require.Equal(t, symbol, msg.Symbol)
			require.Equal(t, models.Update, msg.Type)
			require.True(t, msg.Last.IsPositive())
			require.True(t, msg.Ask.IsPositive())
			require.True(t, msg.Bid.IsPositive())
			require.True(t, msg.AskSize.IsPositive())
			require.True(t, msg.BidSize.IsPositive())
			require.True(t, msg.Bid.LessThanOrEqual(msg.Ask))
			t.Log("so far so good")
			count++
		default:
		}
	}
}

func TestStream_SubscribeToMarkets(t *testing.T) {

	ftx, ctx, done := prepForTest()

	data, err := ftx.Stream.SubscribeToMarkets(ctx)
	require.NoError(t, err)

	asset1, asset2 := "ETH", "BTC"
	symbol := fmt.Sprintf("%s/%s", asset1, asset2)

	count := 0
	for {
		select {
		case <-done:
			require.True(t, count > 0)
			return
		case msg := <-data:
			if msg.Name == symbol {
				require.Equal(t, symbol, msg.Name)
				require.Equal(t, true, msg.Enabled)
				require.Equal(t, asset2, msg.QuoteCurrency)
				require.Equal(t, asset1, msg.BaseCurrency)
				t.Log("so far so good")
			}
			count++
		default:
		}
	}
}

func TestStream_ListMarkets(t *testing.T) {

	ftx, ctx, done := prepForTest()

	data, _ := ftx.Stream.SubscribeToMarkets(ctx)

	count := 0
	for {
		select {
		case <-done:
			require.True(t, count > 0)
			return
		case msg := <-data:
			t.Logf("Name: %s\n", msg.Name)
			t.Logf("Type: %s\n", msg.Type)
			t.Logf("Base Currency:  %s\n", msg.BaseCurrency)
			t.Logf("Quote Currency: %s\n", msg.QuoteCurrency)
			t.Logf("Underlying:     %s\n", msg.Underlying)
			count++
		}
	}
}
func TestStream_SubscribeToTrades(t *testing.T) {

	ftx, ctx, done := prepForTest()

	symbol := "BTC-PERP"
	data, err := ftx.Stream.SubscribeToTrades(ctx, symbol)
	require.NoError(t, err)

	lastID := int64(0)
	for {
		select {
		case <-done:
			return
		case msg := <-data:
			require.Equal(t, symbol, msg.Symbol)
			require.True(t, msg.Price.IsPositive())
			require.True(t, msg.Size.IsPositive())
			require.True(t, msg.ID > lastID)
			lastID = msg.ID
			t.Log("so far so good")
		default:
		}
	}
}

func TestStream_SubscribeToOrderBooks(t *testing.T) {

	ftx, ctx, done := prepForTest()

	symbol := "ETH/BTC"
	data, err := ftx.Stream.SubscribeToOrderBooks(ctx, symbol)
	require.NoError(t, err)

	count := 0
	for {
		select {
		case <-done:
			require.True(t, count > 0)
			return
		case msg := <-data:
			require.Equal(t, symbol, msg.Symbol)
			require.True(t, msg.Type == models.Update || msg.Type == models.Partial)
			require.True(t, len(msg.Bids) > 0 || len(msg.Asks) > 0)
			t.Log("so far so good")
			count++
		default:
		}
	}
}
