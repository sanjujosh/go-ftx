package test

import (
	"os"
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"

	"github.com/uscott/go-ftx/api"
	"github.com/uscott/go-ftx/models"
)

func TestFills_GetFills(t *testing.T) {

	ftx := api.New(
		api.WithAuth(os.Getenv("FTX_PROD_MAIN_KEY"), os.Getenv("FTX_PROD_MAIN_SECRET")),
	)
	err := ftx.SetServerTimeDiff()
	if err != nil {
		t.Fatal(errors.WithStack(err))
	}
	require.NoError(t, err)

	for _, side := range []models.Side{models.Buy, models.Sell} {
		_, err := ftx.Orders.PlaceOrder(&models.OrderParams{
			Market: api.PtrString("BTC-PERP"),
			Side:   api.PtrString(string(side)),
			Type:   api.PtrString(string(models.MarketOrder)),
			Size:   api.PtrDecimal(0.001),
		})
		if err != nil {
			t.Fatal(errors.WithStack(err))
		}
	}
	fills, err := ftx.Fills.GetFills(&models.FillParams{
		Limit: api.PtrInt(10),
	})
	if err != nil {
		t.Fatal(errors.WithStack(err))
	}
	for _, f := range fills {
		if f == nil {
			t.Fatal("nil pointer")
		}
		t.Logf("Fill: %+v\n", *f)
	}
}
