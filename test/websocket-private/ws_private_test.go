package test

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
	"github.com/uscott/go-ftx/api"
	"github.com/uscott/go-ftx/models"
)

const (
	sleepDuration time.Duration = 15 * time.Second
	swap                        = "BTC-PERP"
)

func prepForTest() (*api.Client, *context.Context, chan struct{}) {

	ftx := api.New(
		api.WithAuth(os.Getenv("FTX_PROD_MAIN_KEY"), os.Getenv("FTX_PROD_MAIN_SECRET")),
	)
	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan struct{})
	go func() {
		time.Sleep(sleepDuration)
		cancel()
		done <- struct{}{}
	}()
	return ftx, &ctx, done
}

func TestStream_SubscribeToOrdersAndFills(t *testing.T) {

	ftx, ctx, done := prepForTest()
	defer ftx.CancelAllOrders(&models.CancelAllParams{})

	err := ftx.SetServerTimeDiff()
	if err != nil {
		t.Fatal(errors.WithStack(err))
	}
	require.NoError(t, err)

	go func() {

		time.Sleep(sleepDuration / 2)

		perp, err := ftx.Futures.GetFutureByName(swap)
		if err != nil {
			t.Fatal(errors.WithStack(err))
		}

		bid, _ := perp.Bid.Float64()
		ask, _ := perp.Ask.Float64()

		o, err := ftx.Orders.PlaceOrder(&models.OrderParams{
			Market:   api.PtrString(swap),
			Side:     api.PtrString(string(models.Buy)),
			Price:    api.PtrDecimal(bid - 2),
			Type:     api.PtrString(string(models.LimitOrder)),
			Size:     api.PtrDecimal(0.001),
			PostOnly: api.PtrBool(true),
		})
		if err != nil {
			t.Fatal(errors.WithStack(err))
		}
		oidbid := o.ID

		o, err = ftx.Orders.PlaceOrder(&models.OrderParams{
			Market:   api.PtrString(swap),
			Side:     api.PtrString(string(models.Sell)),
			Price:    api.PtrDecimal(ask + 2),
			Type:     api.PtrString(string(models.LimitOrder)),
			Size:     api.PtrDecimal(0.001),
			PostOnly: api.PtrBool(true),
		})
		if err != nil {
			t.Fatal(errors.WithStack(err))
		}
		oidask := o.ID

		time.Sleep(time.Second)

		_, err = ftx.Orders.ModifyOrder(
			oidbid,
			&models.ModifyOrderParams{
				Price: api.PtrDecimal(bid - 1),
			},
		)
		if err != nil {
			t.Fatal(errors.WithStack(err))
		}
		_, err = ftx.Orders.ModifyOrder(
			oidask,
			&models.ModifyOrderParams{
				Price: api.PtrDecimal(ask + 1),
			},
		)
	}()

	if err != nil {
		t.Fatal(errors.WithStack(err))
	}

	filldata, err := ftx.Stream.SubscribeToFills(*ctx)
	if err != nil {
		t.Fatal(errors.WithStack(err))
	}
	orderdata, err := ftx.Stream.SubscribeToOrders(*ctx)
	if err != nil {
		t.Fatal(errors.WithStack(err))
	}

	for {
		select {
		case <-done:
			return
		case fill, ok := <-filldata:
			if ok {
				t.Logf("Fill: %+v\n", *fill)
			}
		case order, ok := <-orderdata:
			if ok {
				t.Logf("Order: %+v\n", *order)
			}
		default:
			t.Log("waiting ...")
			time.Sleep(100 * time.Millisecond)
		}
	}
}
