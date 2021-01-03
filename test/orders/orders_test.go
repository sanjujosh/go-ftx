package test

import (
	"os"
	"testing"

	"github.com/joho/godotenv"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/uscott/go-ftx/api"
	"github.com/uscott/go-ftx/models"
)

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

	market := "ETH/BTC"

	orders, err := ftx.Orders.GetOrdersHistory(&models.OrdersHistoryParams{
		Market:    &market,
		Limit:     nil,
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

func TestOrders_GetOrderTriggers(t *testing.T) {
	godotenv.Load()

	ftx := api.New(
		api.WithAuth(os.Getenv("FTX_PROD_MAIN_KEY"), os.Getenv("FTX_PROD_MAIN_SECRET")),
	)
	err := ftx.SetServerTimeDiff()
	require.NoError(t, err)

	orderID := int64(1111)

	triggers, err := ftx.Orders.GetOrderTriggers(orderID)

	// 400 - Bad Request, orderID doesn't exist
	assert.Error(t, err)
	assert.Nil(t, triggers)
}

func TestOrders_PlaceOrderModifyAndCancel(t *testing.T) {

	ftx := api.New(
		api.WithAuth(os.Getenv("FTX_PROD_MAIN_KEY"), os.Getenv("FTX_PROD_MAIN_SECRET")),
	)
	err := ftx.SetServerTimeDiff()
	require.NoError(t, err)

	order, err := ftx.Orders.PlaceOrder(&models.OrderParams{
		Market:   api.PtrString("BTC-PERP"),
		Side:     api.PtrString("buy"),
		Price:    api.PtrDecimal(30e3),
		Type:     api.PtrString("limit"),
		Size:     api.PtrDecimal(0.001 / 2),
		PostOnly: api.PtrBool(true),
	})
	if err != nil {
		t.Fatal(errors.WithStack(err))
	}
	t.Logf("Place Order Result: %+v\n", *order)
	orderID := order.ID
	order, err = ftx.Orders.ModifyOrder(
		orderID,
		&models.ModifyOrderParams{
			Price: api.PtrDecimal(29.5e3),
			Size:  api.PtrDecimal(0.001 / 4),
		})
	if err != nil {
		t.Fatal(errors.WithStack(err))
	}
	t.Logf("Modify Order Result: %+v\n", *order)
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

	order, err := ftx.Orders.PlaceTriggerOrder(&models.TriggerOrderParams{
		Market:       api.PtrString("BTC-PERP"),
		Side:         api.PtrString("sell"),
		Size:         api.PtrDecimal(0.001 / 2),
		Type:         api.PtrString("stop"),
		TriggerPrice: api.PtrDecimal(20e3),
		OrderPrice:   api.PtrDecimal(19.5e3),
	})
	if err != nil {
		t.Fatal(errors.WithStack(err))
	}
	t.Logf("Place Trigger Order Result: %+v\n", *order)
	orderID := order.OrderID
	order, err = ftx.Orders.ModifyTriggerOrder(
		orderID,
		&models.ModifyTriggerOrderParams{
			TriggerPrice: api.PtrDecimal(19e3),
			OrderPrice:   api.PtrDecimal(18e3),
		},
	)
	if err != nil {
		t.Fatal(errors.WithStack(err))
	}
	t.Logf("Modify Trigger Order Result: %+v\n", *order)
	orderID = order.OrderID
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

	order, err := ftx.Orders.PlaceOrder(&models.OrderParams{
		Market:   api.PtrString("BTC-PERP"),
		Side:     api.PtrString("buy"),
		Price:    api.PtrDecimal(30e3),
		Type:     api.PtrString("limit"),
		Size:     api.PtrDecimal(0.001 / 2),
		PostOnly: api.PtrBool(true),
	})
	if err != nil {
		t.Fatal(errors.WithStack(err))
	}
	t.Logf("Place Order Result: %+v\n", *order)

	triggerOrder, err := ftx.Orders.PlaceTriggerOrder(&models.TriggerOrderParams{
		Market:       api.PtrString("BTC-PERP"),
		Side:         api.PtrString("sell"),
		Size:         api.PtrDecimal(0.001 / 2),
		Type:         api.PtrString("stop"),
		TriggerPrice: api.PtrDecimal(20e3),
		OrderPrice:   api.PtrDecimal(19.5e3),
	})
	if err != nil {
		t.Fatal(errors.WithStack(err))
	}
	t.Logf("Place Trigger Order Result: %+v\n", *triggerOrder)

	if err = ftx.Orders.CancelAllOrders(&models.CancelAllParams{}); err != nil {
		t.Fatal(errors.WithStack(err))
	}
}
