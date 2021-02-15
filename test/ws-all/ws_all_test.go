package wsalltest

import (
	"testing"

	"github.com/shopspring/decimal"
	"github.com/uscott/go-ftx/api"
	"github.com/uscott/go-ftx/models"
	"github.com/uscott/go-ftx/test"
)

func Test_WsAll(t *testing.T) {

	ftx, ctx, done := test.PrepForTest()

	defer ftx.CancelAllOrders(&models.CancelAllParams{Market: api.PtrString(test.USDTSWAP)})

	var err error
	go test.PlaceSampleOrders(ftx, t, test.USDTSWAP, decimal.NewFromInt(1), &err)
	if err != nil {
		t.Fatal(err)
	}

	symbols := []string{test.USDTSWAP, "BTC-PERP", "BTC/USD"}

	booksC, err := ftx.Stream.SubscribeToOrderBooks(ctx, symbols...)
	if err != nil {
		t.Fatal(err)
	}
	fillsC, err := ftx.Stream.SubscribeToFills(ctx)
	if err != nil {
		t.Fatal(err)
	}
	marketsC, err := ftx.Stream.SubscribeToMarkets(ctx)
	if err != nil {
		t.Fatal(err)
	}
	ordersC, err := ftx.Stream.SubscribeToOrders(ctx, symbols...)
	tickersC, err := ftx.Stream.SubscribeToTickers(ctx, symbols...)
	if err != nil {
		t.Fatal(err)
	}
	tradesC, err := ftx.Stream.SubscribeToTrades(ctx, symbols...)
	if err != nil {
		t.Fatal(err)
	}

	ibooks, ifills, imarkets, iorders, itickers, itrades := 0, 0, 0, 0, 0, 0
	for {
		select {
		case <-done:
			t.Log("Exiting")
			return
		case book := <-booksC:
			if ibooks < test.N {
				t.Logf("Book: %+v", *book)
				ibooks++
			}
		case fill := <-fillsC:
			if ifills < test.N {
				t.Logf("Fill: %+v", *fill)
				ifills++
			}
		case market := <-marketsC:
			if imarkets < test.N {
				t.Logf("Market: %+v", *market)
				imarkets++
			}
		case order := <-ordersC:
			if iorders < test.N {
				t.Logf("Order: %+v", *order)
				iorders++
			}
		case ticker := <-tickersC:
			if itickers < test.N {
				t.Logf("Ticker: %+v", *ticker)
				itickers++
			}
		case trade := <-tradesC:
			if itrades < test.N {
				t.Logf("Trade: %+v", *trade)
				itrades++
			}
		}
	}
}
