package test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/uscott/go-ftx/api"
	"github.com/uscott/go-ftx/models"
)

const (
	runtime time.Duration = 10 * time.Second
	N                     = 5
)

func prepForTest() (*api.Client, *context.Context, chan struct{}) {
	ftx := api.New()
	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan struct{})
	go func() {
		time.Sleep(runtime)
		cancel()
		done <- struct{}{}
	}()
	return ftx, &ctx, done
}

func sleep() { time.Sleep(time.Millisecond) }

func TestStream_SubscribeToTickers(t *testing.T) {

	ftx, ctx, done := prepForTest()

	symbol := "BTC-PERP"
	data, err := ftx.Stream.SubscribeToTickers(*ctx, symbol)
	require.NoError(t, err)

	count := 0
	for {
		select {
		case <-done:
			require.True(t, count > 0)
			return
		case msg := <-data:
			if count > N {
				return
			}
			require.Equal(t, symbol, msg.Symbol)
			require.Equal(t, models.Update, msg.ResponseType)
			require.True(t, msg.Last.IsPositive())
			require.True(t, msg.Ask.IsPositive())
			require.True(t, msg.Bid.IsPositive())
			require.True(t, msg.AskSize.IsPositive())
			require.True(t, msg.BidSize.IsPositive())
			require.True(t, msg.Bid.LessThanOrEqual(msg.Ask))
			t.Logf("Msg: %+v\n", *msg)
			count++
		default:
			sleep()
		}
	}
}

func TestStream_SubscribeToMarkets(t *testing.T) {

	ftx, ctx, done := prepForTest()

	data, err := ftx.Stream.SubscribeToMarkets(*ctx)
	require.NoError(t, err)

	asset1, asset2 := "BTC", "USD"
	symbol := fmt.Sprintf("%s/%s", asset1, asset2)

	count := 0
	for {
		select {
		case <-done:
			require.True(t, count > 0)
			return
		case msg := <-data:
			if msg.Name == symbol {
				if count > N {
					return
				}
				require.Equal(t, symbol, msg.Name)
				require.Equal(t, true, msg.Enabled)
				require.Equal(t, asset2, msg.QuoteCurrency)
				require.Equal(t, asset1, msg.BaseCurrency)
				t.Logf("Msg: %+v\n", *msg)
				count++
			}
		default:
			sleep()
		}
	}
}

func TestStream_ListMarkets(t *testing.T) {

	ftx, ctx, done := prepForTest()

	data, _ := ftx.Stream.SubscribeToMarkets(*ctx)

	count := 0
	for {
		select {
		case <-done:
			require.True(t, count > 0)
			return
		case msg := <-data:
			if count > N {
				return
			}
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
	data, err := ftx.Stream.SubscribeToTrades(*ctx, symbol)
	require.NoError(t, err)

	count, lastID := 0, int64(0)
	for {
		select {
		case <-done:
			return
		case msg := <-data:
			if count > N {
				return
			}
			require.Equal(t, symbol, msg.Symbol)
			require.True(t, msg.Price.IsPositive())
			require.True(t, msg.Size.IsPositive())
			require.True(t, msg.ID > lastID)
			lastID = msg.ID
			t.Logf("Msg: %+v\n", *msg)
		default:
			sleep()
		}
	}
}

func TestStream_SubscribeToOrderBooks(t *testing.T) {

	ftx, ctx, done := prepForTest()

	symbol := "BTC-PERP"
	data, err := ftx.Stream.SubscribeToOrderBooks(*ctx, symbol)
	require.NoError(t, err)

	count := 0
	for {
		select {
		case <-done:
			require.True(t, count > 0)
			return
		case msg := <-data:
			if count > N {
				return
			}

			require.Equal(t, symbol, msg.Symbol)
			require.True(t,
				msg.ResponseType == models.Update || msg.ResponseType == models.Partial)
			require.True(t, len(msg.Bids) > 0 || len(msg.Asks) > 0)
			t.Logf("Msg: %+v\n", *msg)
			count++
		default:
			sleep()
		}
	}
}
