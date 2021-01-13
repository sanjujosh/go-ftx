package test

import (
	"os"
	"testing"

	"github.com/joho/godotenv"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/uscott/go-ftx/api"
	"github.com/uscott/go-ftx/models"
)

const swap = "BTC-PERP"

func client(t *testing.T) *api.Client {

	ftx := api.New(
		api.WithAuth(os.Getenv("FTX_PROD_MAIN_KEY"), os.Getenv("FTX_PROD_MAIN_SECRET")),
	)
	if err := ftx.SetServerTimeDiff(); err != nil {
		t.Fatal(errors.WithStack(err))
	}
	return ftx
}

func TestOrders_GetOpenOrders(t *testing.T) {
	godotenv.Load()

	ftx := api.New(
		api.WithAuth(os.Getenv("FTX_PROD_MAIN_KEY"), os.Getenv("FTX_PROD_MAIN_SECRET")),
	)
	err := ftx.SetServerTimeDiff()
	require.NoError(t, err)

	market := "ETH/BTC"

	orders, err := ftx.Orders.GetOpenOrders(market)
	assert.NoError(t, err)
	assert.NotNil(t, orders)
}

func TestOrders_GetOrdersHistory(t *testing.T) {
	godotenv.Load()

	ftx := api.New(
		api.WithAuth(os.Getenv("FTX_PROD_MAIN_KEY"), os.Getenv("FTX_PROD_MAIN_SECRET")),
	)
	err := ftx.SetServerTimeDiff()
	require.NoError(t, err)

	market, limit := swap, 10

	orders, err := ftx.Orders.GetOrdersHistory(&models.OrdersHistoryParams{
		Market:    &market,
		Limit:     &limit,
		StartTime: nil,
		EndTime:   nil,
	})
	assert.NoError(t, err)
	assert.NotNil(t, orders)
}

func TestOrders_GetOpenTriggerOrders(t *testing.T) {
	godotenv.Load()

	ftx := api.New(
		api.WithAuth(os.Getenv("FTX_PROD_MAIN_KEY"), os.Getenv("FTX_PROD_MAIN_SECRET")),
	)
	err := ftx.SetServerTimeDiff()
	require.NoError(t, err)

	market := "ETH/BTC"
	orderType := models.Stop

	orders, err := ftx.Orders.GetOpenTriggerOrders(&models.OpenTriggerOrdersParams{
		Market: &market,
		Type:   &orderType,
	})
	assert.NoError(t, err)
	assert.NotNil(t, orders)
}

func TestOrders_GetTriggerOrderTriggers(t *testing.T) {
	godotenv.Load()

	ftx := api.New(
		api.WithAuth(os.Getenv("FTX_PROD_MAIN_KEY"), os.Getenv("FTX_PROD_MAIN_SECRET")),
	)
	err := ftx.SetServerTimeDiff()
	require.NoError(t, err)

	orderID := int64(1111)

	triggers, err := ftx.Orders.GetTriggerOrderTriggers(orderID)

	// 400 - Bad Request, orderID doesn't exist
	assert.Error(t, err)
	assert.Nil(t, triggers)
}

func TestOrders_GetTriggerOrdersHistory(t *testing.T) {

	params := &models.TriggerOrdersHistoryParams{
		Market: api.PtrString(swap),
		Limit:  api.PtrInt(10),
	}
	hist, err := client(t).GetTriggerOrdersHistory(params)
	if err != nil {
		t.Fatal(errors.WithStack(err))
	}
	for i, h := range hist {
		if i > 9 {
			return
		}
		t.Logf("Trigger Order: %+v\n", *h)
	}
}

func TestOrders_PlaceOrderModifyAndCancel(t *testing.T) {

	ftx := api.New(
		api.WithAuth(os.Getenv("FTX_PROD_MAIN_KEY"), os.Getenv("FTX_PROD_MAIN_SECRET")),
	)
	err := ftx.SetServerTimeDiff()
	require.NoError(t, err)

	future := models.Future{}
	err = ftx.Futures.GetFutureByName(swap, &future)
	if err != nil {
		t.Fatal(errors.WithStack(err))
	}

	bid, _ := future.Bid.Float64()

	order, price := models.Order{}, decimal.NewFromFloat(bid-100)
	orderType, size := models.LimitOrder, decimal.NewFromFloat(0.01)

	err = ftx.Orders.PlaceOrder(&models.OrderParams{
		Market:   api.PtrString(swap),
		Side:     api.PtrString(string(models.Buy)),
		Price:    &price,
		Type:     api.PtrString(string(orderType)),
		Size:     &size,
		PostOnly: api.PtrBool(true),
	}, &order)
	if err != nil {
		t.Fatal(errors.WithStack(err))
	}
	t.Logf("Place Order Result: %+v\n", order)
	orderID := order.ID
	price = price.Sub(decimal.NewFromInt(100))
	err = ftx.Orders.ModifyOrder(
		orderID,
		&models.ModifyOrderParams{
			Price: &price,
		},
		&order,
	)
	if err != nil {
		t.Fatal(errors.WithStack(err))
	}
	t.Logf("Modify Order Result: %+v\n", order)
	orderID = order.ID
	if err = ftx.Orders.CancelOrder(orderID); err != nil {
		t.Fatal(errors.WithStack(err))
	}
}

func TestOrders_PlaceTriggerOrderModifyAndCancel(t *testing.T) {

	ftx := api.New(
		api.WithAuth(os.Getenv("FTX_PROD_MAIN_KEY"), os.Getenv("FTX_PROD_MAIN_SECRET")),
	)
	err := ftx.SetServerTimeDiff()
	require.NoError(t, err)

	future := models.Future{}
	err = ftx.Futures.GetFutureByName(swap, &future)
	if err != nil {
		t.Fatal(errors.WithStack(err))
	}

	bid, _ := future.Bid.Float64()

	order, size := models.TriggerOrder{}, decimal.NewFromFloat(0.01)
	triggerPrice := decimal.NewFromFloat(bid - 5e3)
	orderPrice := triggerPrice.Sub(decimal.NewFromFloat(1e3))
	err = ftx.Orders.PlaceTriggerOrder(&models.TriggerOrderParams{
		Market:       api.PtrString(swap),
		Side:         api.PtrString(string(models.Sell)),
		Size:         &size,
		Type:         api.PtrString(string(models.Stop)),
		TriggerPrice: &triggerPrice,
		OrderPrice:   &orderPrice,
	}, &order)
	if err != nil {
		t.Fatal(errors.WithStack(err))
	}
	t.Logf("Place Trigger Order Result: %+v\n", order)
	orderID := order.ID
	triggerPrice = triggerPrice.Sub(decimal.NewFromInt(100))
	orderPrice = orderPrice.Sub(decimal.NewFromInt(100))
	err = ftx.Orders.ModifyTriggerOrder(
		orderID,
		&models.ModifyTriggerOrderParams{
			TriggerPrice: &triggerPrice,
			OrderPrice:   &orderPrice,
		},
		&order,
	)
	if err != nil {
		t.Fatal(errors.WithStack(err))
	}
	t.Logf("Modify Trigger Order Result: %+v\n", order)
	orderID = order.ID
	if err = ftx.Orders.CancelTriggerOrder(orderID); err != nil {
		t.Fatal(errors.WithStack(err))
	}
}

func TestOrders_CancelAll(t *testing.T) {

	ftx := api.New(
		api.WithAuth(os.Getenv("FTX_PROD_MAIN_KEY"), os.Getenv("FTX_PROD_MAIN_SECRET")),
	)
	err := ftx.SetServerTimeDiff()
	require.NoError(t, err)

	future := models.Future{}
	err = ftx.Futures.GetFutureByName(swap, &future)
	if err != nil {
		t.Fatal(errors.WithStack(err))
	}

	bid, _ := future.Bid.Float64()

	order := models.Order{}
	err = ftx.Orders.PlaceOrder(&models.OrderParams{
		Market:   api.PtrString(swap),
		Side:     api.PtrString(string(models.Buy)),
		Price:    api.PtrDecimal(bid - 100),
		Type:     api.PtrString(string(models.LimitOrder)),
		Size:     api.PtrDecimal(0.01),
		PostOnly: api.PtrBool(true),
	}, &order)
	if err != nil {
		t.Fatal(errors.WithStack(err))
	}
	t.Logf("Place Order Result: %+v\n", order)

	triggerOrder := models.TriggerOrder{}
	err = ftx.Orders.PlaceTriggerOrder(&models.TriggerOrderParams{
		Market:       api.PtrString(swap),
		Side:         api.PtrString(string(models.Sell)),
		Size:         api.PtrDecimal(0.01),
		Type:         api.PtrString(string(models.Stop)),
		TriggerPrice: api.PtrDecimal(bid - 5e3),
		OrderPrice:   api.PtrDecimal(bid - 6e3),
	}, &triggerOrder)
	if err != nil {
		t.Fatal(errors.WithStack(err))
	}
	t.Logf("Place Trigger Order Result: %+v\n", triggerOrder)

	err = ftx.Orders.CancelAllOrders(&models.CancelAllParams{
		Market: api.PtrString(swap),
		Side:   api.PtrString("buy"),
	})
	if err != nil {
		t.Fatal(errors.WithStack(err))
	}

	err = ftx.Orders.CancelAllOrders(&models.CancelAllParams{
		Market: api.PtrString(swap),
		Side:   api.PtrString("sell"),
	})
	if err != nil {
		t.Fatal(errors.WithStack(err))
	}
}
