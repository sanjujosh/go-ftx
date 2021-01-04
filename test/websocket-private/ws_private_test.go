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
	ftx := api.New()
	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan struct{})
	go func() {
		time.Sleep(sleepDuration)
		cancel()
		done <- struct{}{}
	}()
	return ftx, &ctx, done
}

func TestStream_SubscribeToFills(t *testing.T) {

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
	err := ftx.SetServerTimeDiff()
	if err != nil {
		t.Fatal(errors.WithStack(err))
	}
	data, err := ftx.Stream.SubscribeToFills(ctx)
	require.NoError(t, err)
	perp, err := ftx.Futures.GetFutureByName(swap)
	if err != nil {
		t.Fatal(errors.WithStack(err))
	}
	ask, _ := perp.Ask.Float64()
	_, err = ftx.Orders.PlaceOrder(&models.OrderParams{
		Market:   api.PtrString(swap),
		Side:     api.PtrString(string(models.Buy)),
		Price:    api.PtrDecimal(ask + 1),
		Type:     api.PtrString(string(models.LimitOrder)),
		Size:     api.PtrDecimal(0.001),
		PostOnly: api.PtrBool(true),
	})
	if err != nil {
		t.Fatal(errors.WithStack(err))
	}
	bid, _ := perp.Bid.Float64()
	_, err = ftx.Orders.PlaceOrder(&models.OrderParams{
		Market:   api.PtrString(swap),
		Side:     api.PtrString(string(models.Sell)),
		Price:    api.PtrDecimal(bid - 1),
		Type:     api.PtrString(string(models.LimitOrder)),
		Size:     api.PtrDecimal(0.001),
		PostOnly: api.PtrBool(true),
	})
	if err != nil {
		t.Fatal(errors.WithStack(err))
	}

	for {
		select {
		case <-done:
			return
		case msg := <-data:
			if msg != nil {
				t.Logf("Data: %+v\n", *msg)
			}
		default:
		}
	}
}
