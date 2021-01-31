package test

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/shopspring/decimal"
	"github.com/uscott/go-ftx/api"
	"github.com/uscott/go-ftx/models"
)

const (
	SleepDuration time.Duration = 10 * time.Second
	USDTSWAP                    = "USDT-PERP"
)

func PrepForTest() (*api.Client, *context.Context, chan struct{}) {

	ftx := api.New(
		api.WithAuth(os.Getenv("FTX_PROD_MAIN_KEY"), os.Getenv("FTX_PROD_MAIN_SECRET")),
	)
	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan struct{})
	go func() {
		time.Sleep(SleepDuration)
		cancel()
		done <- struct{}{}
	}()
	return ftx, &ctx, done
}

func PlaceSampleOrders(
	ftx *api.Client, t *testing.T, future string, size decimal.Decimal, err *error) {

	if err == nil {
		panic("Nil pointer")
	}

	time.Sleep(SleepDuration / 2)

	perp := &models.Future{}
	*err = ftx.Futures.GetFutureByName(USDTSWAP, perp)
	if *err != nil {
		return
	}

	bid := perp.Bid
	ask := perp.Ask
	incr := perp.PriceIncrement
	o := &models.Order{}

	e := ftx.Orders.PlaceOrder(&models.OrderParams{
		Market:   api.PtrString(USDTSWAP),
		Side:     api.PtrString(string(models.Buy)),
		Price:    api.PtrDecimal(bid.Sub(incr)),
		Type:     api.PtrString(string(models.LimitOrder)),
		Size:     &size,
		PostOnly: api.PtrBool(false),
	}, o)
	if e != nil {
		t.Log(e.Error())
	}
	oidbid := o.ID

	e = ftx.Orders.PlaceOrder(&models.OrderParams{
		Market:   api.PtrString(USDTSWAP),
		Side:     api.PtrString(string(models.Sell)),
		Price:    api.PtrDecimal(ask.Add(incr)),
		Type:     api.PtrString(string(models.LimitOrder)),
		Size:     &size,
		PostOnly: api.PtrBool(false),
	}, o)
	if e != nil {
		t.Log(e.Error())
	}
	oidask := o.ID

	time.Sleep(time.Second)

	e = ftx.Orders.ModifyOrder(
		oidbid,
		&models.ModifyOrderParams{
			Price: api.PtrDecimal(ask.Add(incr)),
			Size:  &size,
		},
		o,
	)
	if e != nil {
		t.Log(e.Error())
	}
	e = ftx.Orders.ModifyOrder(
		oidask,
		&models.ModifyOrderParams{
			Price: api.PtrDecimal(bid.Sub(incr)),
			Size:  &size,
		},
		o,
	)
	if e != nil {
		t.Log(e.Error())
	}
}
