package testorders

import (
	"os"
	"path"
	"testing"
	"time"

	"github.com/joho/godotenv"
	"github.com/pkg/errors"
	"github.com/sanjujosh/go-ftx/api"
	"github.com/sanjujosh/go-ftx/models"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

const (
	swap = "USDT-PERP"
	N    = 5
)

func client(t *testing.T) *api.Client {

	// Load the config from user's home directory
	homeDir, err := os.UserHomeDir()
	if err != nil {
		t.Fatal(err)
	}

	err = godotenv.Load(path.Join(homeDir, ".custom_project_config", ".go-ftx", ".env"))
	if err != nil {
		t.Fatal(err)
	}

	ftx := api.New(
		api.WithAuth(os.Getenv("FTX_PROD_MAIN_KEY"), os.Getenv("FTX_PROD_MAIN_SECRET")),
		api.SetSubAccount(os.Getenv("FTX_PROD_MAIN_ACC")),
	)
	if err := ftx.SetServerTimeDiff(); err != nil {
		t.Fatal(errors.WithStack(err))
	}
	return ftx
}

func TestOrders_GetOpenOrders(t *testing.T) {

	ftx := client(t)

	err := ftx.SetServerTimeDiff()
	if err != nil {
		t.Fatal(err)
	}

	orders, err := ftx.Orders.GetOpenOrders("")
	if err != nil {
		t.Fatal(err)
	}
	if orders == nil {
		t.Fatal("Orders should not be nil")
	}

	for i, o := range orders {
		if i > N {
			return
		}
		t.Logf("Order: %+v\n", *o)
	}
}

func TestOrders_GetOrdersHistory(t *testing.T) {

	ftx := api.New(
		api.WithAuth(os.Getenv("FTX_PROD_MAIN_KEY"), os.Getenv("FTX_PROD_MAIN_SECRET")),
	)
	err := ftx.SetServerTimeDiff()
	if err != nil {
		t.Fatal(err)
	}

	limit := 10

	orders, err := ftx.Orders.GetOrdersHistory(&models.OrdersHistoryParams{
		Market:    nil,
		Limit:     &limit,
		StartTime: nil,
		EndTime:   nil,
	})
	if err != nil {
		t.Fatal(err)
	}
	if orders == nil {
		t.Fatal("Orders should not be nil")
	}

	for i, o := range orders {
		if i > N {
			return
		}
		t.Logf("Order: %+v\n", *o)
	}
}

func TestOrders_GetOpenTriggerOrders(t *testing.T) {

	ftx := api.New(
		api.WithAuth(os.Getenv("FTX_PROD_MAIN_KEY"), os.Getenv("FTX_PROD_MAIN_SECRET")),
	)
	err := ftx.SetServerTimeDiff()
	if err != nil {
		t.Fatal(err)
	}

	orders, err := ftx.Orders.GetOpenTriggerOrders(nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	if orders == nil {
		t.Fatal("Orders should not be nil")
	}

	for i, o := range orders {
		if i > N {
			return
		}
		t.Logf("Order: %+v\n", *o)
	}
}

func TestOrders_GetTriggerOrderTriggers(t *testing.T) {

	ftx := api.New(
		api.WithAuth(os.Getenv("FTX_PROD_MAIN_KEY"), os.Getenv("FTX_PROD_MAIN_SECRET")),
	)
	err := ftx.SetServerTimeDiff()
	if err != nil {
		t.Fatal(err)
	}

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
		if i > N {
			return
		}
		t.Logf("Trigger Order: %+v\n", *h)
	}
}

func TestOrders_PlaceOrderModifyAndCancel(t *testing.T) {

	ftx := client(t)

	err := ftx.SetServerTimeDiff()
	if err != nil {
		t.Fatal(err)
	}

	future := models.Future{}
	err = ftx.Futures.GetFutureByName(swap, &future)
	if err != nil {
		t.Fatal(errors.WithStack(err))
	}

	bid, _ := future.Bid.Float64()

	order, price := models.Order{}, decimal.NewFromFloat(bid-0.50)
	orderType, size := models.LimitOrder, decimal.NewFromFloat(1)

	err = ftx.Orders.PlaceOrder(&models.OrderParams{
		Market:   swap,
		Side:     models.Buy,
		Price:    price,
		Type:     orderType,
		Size:     size,
		PostOnly: true,
	}, &order)
	if err != nil {
		t.Fatal(errors.WithStack(err))
	}
	t.Logf("Place Order Result: %+v\n", order)
	orderID := order.ID
	price = price.Sub(decimal.NewFromFloat(0.25))
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
	success, err := ftx.Orders.CancelOrder(orderID)
	if err != nil {
		t.Fatal(errors.WithStack(err))
	}
	t.Logf("Cancel Result: %+v\n", success)
}

func TestOrders_PlaceTriggerOrderModifyAndCancel(t *testing.T) {

	ftx := api.New(
		api.WithAuth(os.Getenv("FTX_PROD_MAIN_KEY"), os.Getenv("FTX_PROD_MAIN_SECRET")),
	)
	err := ftx.SetServerTimeDiff()
	if err != nil {
		t.Fatal(err)
	}

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
	success, err := ftx.Orders.CancelTriggerOrder(orderID)
	if err != nil {
		t.Fatal(errors.WithStack(err))
	}
	t.Logf("Cancel Result: %+v\n", success)
}

func TestOrders_CancelAll(t *testing.T) {

	ftx := api.New(
		api.WithAuth(os.Getenv("FTX_PROD_MAIN_KEY"), os.Getenv("FTX_PROD_MAIN_SECRET")),
	)
	err := ftx.SetServerTimeDiff()
	if err != nil {
		t.Fatal(err)
	}

	future, order1, order2 := models.Future{}, models.Order{}, models.Order{}
	contract := "BTC-0924"

	for c, o := range map[string]*models.Order{swap: &order1, contract: &order2} {

		if err = ftx.Futures.GetFutureByName(c, &future); err != nil {
			t.Fatal(errors.WithStack(err))
		}

		price := future.Bid.Sub(decimal.NewFromFloat(1000))

		err = ftx.Orders.PlaceOrder(&models.OrderParams{
			Market:   c,
			Side:     models.Buy,
			Price:    price,
			Type:     models.LimitOrder,
			Size:     decimal.NewFromFloat(0.01),
			PostOnly: true,
		}, o)

		if err != nil {
			t.Fatal(errors.WithStack(err))
		}

		t.Logf("\nPlace %s Order Result: %+v\n", c, *o)
		time.Sleep(time.Second)
	}

	success, err := ftx.Orders.CancelAllOrders(&models.CancelAllParams{
		Market: api.PtrString(swap),
	})

	if err != nil {
		t.Fatal(errors.WithStack(err))
	}

	t.Logf("\nCancel All Orders %s Result: %+v\n", swap, success)

	for c, o := range map[string]*models.Order{swap: &order1, contract: &order2} {

		if err = ftx.Orders.GetOrderStatus(o.ID, o); err != nil {
			t.Fatal(err)
		}

		t.Logf("\n%s Order Status: %+v\n", c, *o)
		time.Sleep(time.Second)
	}

	success, err = ftx.Orders.CancelAllOrders(&models.CancelAllParams{
		Market: api.PtrString(contract),
	})

	if err != nil {
		t.Fatal(errors.WithStack(err))
	}

	t.Logf("\nCancel All Orders %s Result: %+v\n", contract, success)
	time.Sleep(time.Second)

	if err = ftx.Orders.GetOrderStatus(order2.ID, &order2); err != nil {
		t.Fatal(err)
	}

	t.Logf("\n%s Order Status: %+v\n", contract, order2)
}
