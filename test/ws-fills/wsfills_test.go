package testwsfills

import (
	"testing"
	"time"

	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/require"
	"github.com/uscott/go-ftx/api"
	"github.com/uscott/go-ftx/models"
	"github.com/uscott/go-ftx/test"
)

func Test_Fills(t *testing.T) {

	ftx, ctx, done := test.PrepForTest()
	defer ftx.CancelAllOrders(&models.CancelAllParams{Market: api.PtrString(test.USDTSWAP)})

	err := ftx.SetServerTimeDiff()
	if err != nil {
		t.Fatal(errors.WithStack(err))
	}
	require.NoError(t, err)

	go test.PlaceSampleOrders(ftx, t, test.USDTSWAP, decimal.NewFromInt(1), &err)

	if err != nil {
		t.Fatal(err)
	}

	filldata, err := ftx.Stream.SubscribeToFills(*ctx)
	if err != nil {
		t.Fatal(errors.WithStack(err))
	}

	for {
		select {

		case <-done:
			return
		case fill := <-filldata:
			t.Log("yes!")
			t.Logf("Fill: %+v\n", *fill)
		default:
			time.Sleep(time.Millisecond)

		}
	}
}
