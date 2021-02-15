package testwsorders

import (
	"testing"
	"time"

	"github.com/shopspring/decimal"
	"github.com/uscott/go-ftx/api"
	"github.com/uscott/go-ftx/models"
	"github.com/uscott/go-ftx/test"
)

func Test_Orders(t *testing.T) {

	ftx, ctx, done := test.PrepForTest()
	defer ftx.CancelAllOrders(&models.CancelAllParams{Market: api.PtrString(test.USDTSWAP)})

	var err error

	go test.PlaceSampleOrders(ftx, t, test.USDTSWAP, decimal.NewFromInt(1), &err)
	if err != nil {
		t.Fatal(err)
	}

	orderC, err := ftx.Stream.SubscribeToOrders(ctx)
	if err != nil {
		t.Fatal(err)
	}

	for {
		select {
		case <-done:
			return
		case order := <-orderC:
			t.Logf("Order: %+v\n", *order)
		default:
			time.Sleep(time.Millisecond)
		}
	}
}
