package models

import (
	"time"

	"github.com/shopspring/decimal"
)

type FundingPayment struct {
	Future  string          `json:"future"`
	ID      int64           `json:"id"`
	Payment decimal.Decimal `json:"payment"`
	Time    time.Time       `json:"time"`
	Rate    decimal.Decimal `json:"rate"`
}
