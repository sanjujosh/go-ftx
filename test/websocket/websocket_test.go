package testws

import (
	"context"
	"testing"
	"time"

	"github.com/uscott/go-ftx/api"
	"github.com/uscott/go-ftx/test"
)

func sleep() { time.Sleep(time.Microsecond) }

func Test_WS(t *testing.T) {

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	client := api.New()
	client.Stream.SetDebugMode(true)

	symbols := []string{"BTC-PERP", "BTC/USD", "FTT-PERP", "FTT/USD"}

	tickerC, err := client.Stream.SubscribeToTickers(ctx, symbols...)
	if err != nil {
		t.Fatal(err)
	}

	marketsC, err := client.Stream.SubscribeToMarkets(ctx)
	if err != nil {
		t.Fatal(err)
	}

	tradesC, err := client.Stream.SubscribeToTrades(ctx, symbols...)
	if err != nil {
		t.Fatal(err)
	}

	booksC, err := client.Stream.SubscribeToOrderBooks(ctx, symbols...)
	if err != nil {
		t.Fatal(err)
	}

	done := test.MakeDoneChan()

	for {
		select {
		case <-done:
			t.Log("Exiting")
			return
		case ticker := <-tickerC:
			if ticker != nil {
				t.Logf("Ticker: %+v", *ticker)
			}
		case market := <-marketsC:
			if market != nil {
				t.Logf("Market: %+v", *market)
			}
		case trade := <-tradesC:
			if trade != nil {
				t.Logf("Trade: %+v", *trade)
			}
		case book := <-booksC:
			if book != nil {
				t.Logf("Book: %+v", *book)
			}
		default:
			sleep()
		}
	}
}
