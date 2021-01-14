package models

import (
	"time"

	"github.com/shopspring/decimal"
)

type FillParams struct {
	Market    *string `json:"market,omitempty"`
	Limit     *int    `json:"limit,omitempty"`
	StartTime *int64  `json:"start_time,omitempty"`
	EndTime   *int64  `json:"end_time,omitempty"`
	Order     *string `json:"order,omitempty"`
	OrderID   *int64  `json:"orderId,omitempty"`
}

type Fill struct {
	Fee           float64         `json:"fee"`
	FeeCurrency   string          `json:"feeCurrency"`
	FeeRate       float64         `json:"feeRate"`
	Future        string          `json:"future"`
	ID            int64           `json:"id"`
	Liquidity     string          `json:"liquidity"`
	Market        string          `json:"market"`
	BaseCurrency  string          `json:"baseCurrency"`
	QuoteCurrency string          `json:"quoteCurrency"`
	OrderID       int64           `json:"orderId"`
	TradeID       int64           `json:"tradeId"`
	Price         decimal.Decimal `json:"price"`
	Side          string          `json:"side"`
	Size          decimal.Decimal `json:"size"`
	Time          time.Time       `json:"time"`
	Type          string          `json:"type"`
}
