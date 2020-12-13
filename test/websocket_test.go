package test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/uscott/goftx"
	"github.com/uscott/goftx/models"
)

func TestStream_SubscribeToTickers(t *testing.T) {

	ftx := goftx.New()

	ctx, cancel := context.WithCancel(context.Background())

	symbol := "ETH/BTC"
	data, err := ftx.Stream.SubscribeToTickers(ctx, symbol)
	require.NoError(t, err)

	done := make(chan struct{})
	go func() {
		time.Sleep(10 * time.Second)
		cancel()
		done <- struct{}{}
	}()
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
			count++
		}
	}
}

func TestStream_SubscribeToMarkets(t *testing.T) {

	ftx := goftx.New()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	data, err := ftx.Stream.SubscribeToMarkets(ctx)
	require.NoError(t, err)

	asset1, asset2 := "ETH", "BTC"
	symbol := fmt.Sprintf("%s/%s", asset1, asset2)

	count := 0
	for msg := range data {
		if msg.Name != symbol {
			continue
		}
		require.Equal(t, symbol, msg.Name)
		require.Equal(t, true, msg.Enabled)
		require.Equal(t, asset2, msg.QuoteCurrency)
		require.Equal(t, asset1, msg.BaseCurrency)
		count++
	}
	require.True(t, count > 0)
}

func TestStream_SubscribeToTrades(t *testing.T) {

	ftx := goftx.New()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	symbol := "BTC-PERP"
	data, err := ftx.Stream.SubscribeToTrades(ctx, symbol)
	require.NoError(t, err)

	go func() {
		<-ctx.Done()
		<-time.After(time.Second)
		if _, ok := <-data; !ok {
			t.Fail()
		}
	}()

	lastID := int64(0)
	for msg := range data {
		require.Equal(t, symbol, msg.Symbol)
		require.True(t, msg.Price.IsPositive())
		require.True(t, msg.Size.IsPositive())
		require.True(t, msg.ID > lastID)
		lastID = msg.ID
	}
}

func TestStream_SubscribeToOrderBooks(t *testing.T) {

	ftx := goftx.New()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	symbol := "ETH/BTC"
	data, err := ftx.Stream.SubscribeToOrderBooks(ctx, symbol)
	require.NoError(t, err)

	go func() {
		<-ctx.Done()
		<-time.After(time.Second)
		if _, ok := <-data; !ok {
			t.Fail()
		}
	}()

	count := 0
	for msg := range data {
		require.Equal(t, symbol, msg.Symbol)
		require.True(t, msg.Type == models.Update || msg.Type == models.Partial)
		require.True(t, len(msg.Bids) > 0 || len(msg.Asks) > 0)
		count++
	}
	require.True(t, count > 0)
}
