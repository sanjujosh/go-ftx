package testacct

import (
	"os"
	"testing"

	"github.com/joho/godotenv"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/uscott/go-ftx/api"
	"github.com/uscott/go-ftx/models"
)

func TestAccount_GetAccountInformation(t *testing.T) {
	godotenv.Load()

	ftx := api.New(
		api.WithAuth(os.Getenv("FTX_PROD_MAIN_KEY"), os.Getenv("FTX_PROD_MAIN_SECRET")),
	)
	err := ftx.SetServerTimeDiff()
	require.NoError(t, err)

	account := models.AccountInformation{}
	err = ftx.Account.GetAccountInformation(&account)
	assert.NoError(t, err)
	assert.NotNil(t, &account)
	t.Logf("Account: %+v\n", account)
}

func TestAccount_GetPositions(t *testing.T) {
	godotenv.Load()

	ftx := api.New(
		api.WithAuth(os.Getenv("FTX_PROD_MAIN_KEY"), os.Getenv("FTX_PROD_MAIN_SECRET")),
	)
	err := ftx.SetServerTimeDiff()
	require.NoError(t, err)

	positions, err := ftx.Account.GetPositions()
	assert.NoError(t, err)
	assert.NotNil(t, positions)
	for _, p := range positions {
		t.Logf("Position: %+v\n", *p)
	}
}

func TestAccount_ChangeAccountLeverage(t *testing.T) {
	godotenv.Load()

	ftx := api.New(
		api.WithAuth(os.Getenv("FTX_PROD_MAIN_KEY"), os.Getenv("FTX_PROD_MAIN_SECRET")),
	)
	err := ftx.SetServerTimeDiff()
	require.NoError(t, err)

	leverage := 10.0

	result, err := ftx.Account.ChangeAccountLeverage(leverage)
	assert.NoError(t, err)
	t.Logf("Result: %s\n", result)

	account := models.AccountInformation{}
	err = ftx.Account.GetAccountInformation(&account)

	assert.NoError(t, err)
	assert.NotNil(t, &account)
	assert.True(t, decimal.NewFromFloat(leverage).Equal(account.Leverage))
}
