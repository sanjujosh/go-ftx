package test

import (
	"os"
	"testing"

	"github.com/joho/godotenv"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/uscott/go-ftx"
)

func TestAccount_GetAccountInformation(t *testing.T) {
	godotenv.Load()

	ftx := ftx.New(
		ftx.WithAuth(os.Getenv("FTX_PROD_MAIN_KEY"), os.Getenv("FTX_PROD_MAIN_SECRET")),
	)
	err := ftx.SetServerTimeDiff()
	require.NoError(t, err)

	account, err := ftx.Account.GetAccountInformation()
	assert.NoError(t, err)
	assert.NotNil(t, account)
}

func TestAccount_GetPositions(t *testing.T) {
	godotenv.Load()

	ftx := ftx.New(
		ftx.WithAuth(os.Getenv("FTX_PROD_MAIN_KEY"), os.Getenv("FTX_PROD_MAIN_SECRET")),
	)
	err := ftx.SetServerTimeDiff()
	require.NoError(t, err)

	positions, err := ftx.Account.GetPositions()
	assert.NoError(t, err)
	assert.NotNil(t, positions)
}

func TestAccount_ChangeAccountLeverage(t *testing.T) {
	godotenv.Load()

	ftx := ftx.New(
		ftx.WithAuth(os.Getenv("FTX_PROD_MAIN_KEY"), os.Getenv("FTX_PROD_MAIN_SECRET")),
	)
	err := ftx.SetServerTimeDiff()
	require.NoError(t, err)

	leverage := decimal.New(10, 0)

	err = ftx.Account.ChangeAccountLeverage(leverage)
	assert.NoError(t, err)

	account, err := ftx.Account.GetAccountInformation()
	assert.NoError(t, err)
	assert.NotNil(t, account)
	assert.True(t, leverage.Equal(account.Leverage))
}
