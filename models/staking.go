package models

import (
	"time"

	"github.com/shopspring/decimal"
)

type Stake struct {
	Coin      string          `json:"coin"`
	CreatedAt time.Time       `json:"createdAt"`
	ID        int64           `json:"id"`
	Size      decimal.Decimal `json:"size"`
}

type UnstakeRequest struct {
	Stake
	Status   UnstakeRequestStatus `json:"status"`
	UnlockAt time.Time            `json:"unlockAt"`
}

type StakeBalance struct {
	Coin               string          `json:"coin"`
	LifetimeRewards    float64         `json:"lifetimeRewards"`
	ScheduledToUnstake decimal.Decimal `json:"scheduledToUnstake"`
	Staked             decimal.Decimal `json:"staked"`
}

type StakingReward struct {
	Coin   string          `json:"coid"`
	ID     int64           `json:"id"`
	Size   decimal.Decimal `json:"size"`
	Status string          `json:"status"`
	Time   time.Time       `json:"time"`
}

type StakeRequestParams struct {
	Coin *string          `json:"coin"`
	Size *decimal.Decimal `json:"size"`
}

type StakeRequest Stake
type UnstakeRequestParams StakeRequestParams
