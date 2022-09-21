package testwsfills

import (
	"testing"
	"time"

	"github.com/sanjujosh/go-ftx/api"
	"github.com/sanjujosh/go-ftx/models"
	"github.com/sanjujosh/go-ftx/test"
	"github.com/shopspring/decimal"
)

func Test_Fills(t *testing.T) {

	ftx, ctx, done := test.PrepForTest()
	defer ftx.CancelAllOrders(&models.CancelAllParams{Market: api.PtrString(test.USDTSWAP)})

	var err error

	go test.PlaceSampleOrders(ftx, t, test.USDTSWAP, decimal.NewFromInt(1), &err)
	if err != nil {
		t.Fatal(err)
	}

	fillC, err := ftx.Stream.SubscribeToFills(ctx)
	if err != nil {
		t.Fatal(err)
	}

	for {
		select {
		case <-done:
			return
		case fill := <-fillC:
			t.Logf("Fill: %+v\n", *fill)
		default:
			time.Sleep(time.Millisecond)
		}
	}
}
