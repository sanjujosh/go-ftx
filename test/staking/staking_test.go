package test

import (
	"os"
	"testing"

	"github.com/pkg/errors"
	"github.com/uscott/go-ftx/api"
)

func client(t *testing.T) *api.Client {

	ftx := api.New(
		api.WithAuth(os.Getenv("FTX_PROD_MAIN_KEY"), os.Getenv("FTX_PROD_MAIN_SECRET")),
	)
	if err := ftx.SetServerTimeDiff(); err != nil {
		t.Fatal(errors.WithStack(err))
	}
	return ftx
}

func TestStaking_GetStakes(t *testing.T) {

	ftx := client(t)

	stakes, err := ftx.Staking.GetStakes()
	if err != nil {
		t.Fatal(errors.WithStack(err))
	}
	for i, s := range stakes {
		if i > 9 {
			return
		}
		t.Logf("Stake: %+v\n", *s)
	}
}

func TestStaking_GetUnstakeRequests(t *testing.T) {

	requests, err := client(t).Staking.GetUnstakeRequests()
	if err != nil {
		t.Fatal(errors.WithStack(err))
	}
	for i, r := range requests {
		if i > 9 {
			return
		}
		t.Logf("Request: %+v\n", *r)
	}
}

func TestStaking_GetStakeBalances(t *testing.T) {

	balances, err := client(t).Staking.GetStakeBalances()
	if err != nil {
		t.Fatal(errors.WithStack(err))
	}
	for i, b := range balances {
		if i > 9 {
			return
		}
		t.Logf("Balance: %+v\n", *b)
	}
}

func TestStaking_RequestUnstake(t *testing.T) {
}

func TestStaking_CancelUnstakeRequest(t *testing.T) {
}

func TestStaking_GetStakingRewards(t *testing.T) {

	rewards, err := client(t).Staking.GetStakingRewards()
	if err != nil {
		t.Fatal(errors.WithStack(err))
	}
	for i, r := range rewards {
		if i > 9 {
			return
		}
		t.Logf("Balance: %+v\n", *r)
	}
}
