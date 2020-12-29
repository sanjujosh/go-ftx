package test

import (
	"os"
	"testing"

	"github.com/joho/godotenv"
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
	err := api.SetServerTimeDiff()
	require.NoError(t, err)

	market := "ETH/BTC"

	orders, err := api.Orders.GetOpenOrders(market)
	assert.NoError(t, err)
	assert.NotNil(t, orders)
}

func TestOrders_GetOrdersHistory(t *testing.T) {
	godotenv.Load()

	ftx := api.New(
		api.WithAuth(os.Getenv("FTX_PROD_MAIN_KEY"), os.Getenv("FTX_PROD_MAIN_SECRET")),
	)
	err := api.SetServerTimeDiff()
	require.NoError(t, err)

	market := "ETH/BTC"

	orders, err := api.Orders.GetOrdersHistory(&models.GetOrdersHistoryParams{
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
	err := api.SetServerTimeDiff()
	require.NoError(t, err)

	market := "ETH/BTC"
	orderType := models.Stop

	orders, err := api.Orders.GetOpenTriggerOrders(&models.GetOpenTriggerOrdersParams{
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
	err := api.SetServerTimeDiff()
	require.NoError(t, err)

	orderID := int64(1111)

	triggers, err := api.Orders.GetOrderTriggers(orderID)

	// 400 - Bad Request, orderID doesn't exist
	assert.Error(t, err)
	assert.Nil(t, triggers)
}
