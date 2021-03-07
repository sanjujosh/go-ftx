package models

import (
	"time"

	"github.com/shopspring/decimal"
)

type FundingPaymentParams struct {
	StartTime *int64  `json:"start_time,omitempty"`
	EndTime   *int64  `json:"end_time,omitempty"`
	Future    *string `json:"future,omitempty"`
}

type FundingPayment struct {
	Future  string          `json:"future"`
	ID      int64           `json:"id"`
	Payment decimal.Decimal `json:"payment"`
	Time    time.Time       `json:"time"`
	Rate    decimal.Decimal `json:"rate"`
}
