package test

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/require"
	"github.com/uscott/go-ftx/api"
	"github.com/uscott/go-ftx/models"
)

const (
	sleepDuration time.Duration = 10 * time.Second
	swap                        = "USDT-PERP"
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

		perp := &models.Future{}
		err := ftx.Futures.GetFutureByName(swap, perp)
		if err != nil {
			t.Fatal(errors.WithStack(err))
		}

		bid := perp.Bid
		ask := perp.Ask
		incr := perp.PriceIncrement
		size := decimal.NewFromInt(1)
		o := &models.Order{}

		err = ftx.Orders.PlaceOrder(&models.OrderParams{
			Market:   api.PtrString(swap),
			Side:     api.PtrString(string(models.Buy)),
			Price:    api.PtrDecimal(bid.Sub(incr)),
			Type:     api.PtrString(string(models.LimitOrder)),
			Size:     &size,
			PostOnly: api.PtrBool(false),
		}, o)
		if err != nil {
			t.Log(err.Error())
		}
		oidbid := o.ID

		err = ftx.Orders.PlaceOrder(&models.OrderParams{
			Market:   api.PtrString(swap),
			Side:     api.PtrString(string(models.Sell)),
			Price:    api.PtrDecimal(ask.Add(incr)),
			Type:     api.PtrString(string(models.LimitOrder)),
			Size:     &size,
			PostOnly: api.PtrBool(false),
		}, o)
		if err != nil {
			t.Log(err.Error())
		}
		oidask := o.ID

		time.Sleep(time.Second)

		err = ftx.Orders.ModifyOrder(
			oidbid,
			&models.ModifyOrderParams{
				Price: api.PtrDecimal(ask.Add(incr)),
				Size:  &size,
			},
			o,
		)
		if err != nil {
			t.Log(err.Error())
		}
		err = ftx.Orders.ModifyOrder(
			oidask,
			&models.ModifyOrderParams{
				Price: api.PtrDecimal(bid.Sub(incr)),
				Size:  &size,
			},
			o,
		)
		if err != nil {
			t.Log(err.Error())
		}
	}()

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
		case fill := <-filldata:
			t.Log("yes!")
			t.Logf("Fill: %+v\n", *fill)
		case order := <-orderdata:
			t.Log("yay!")
			t.Logf("Order: %+v\n", *order)
		default:
			time.Sleep(time.Millisecond)
		}
	}
}
