package models

import (
	"time"

	"github.com/shopspring/decimal"
)

type Order struct {
	ID            int64           `json:"id"`
	Market        string          `json:"market"`
	Type          OrderType       `json:"type"`
	Side          Side            `json:"side"`
	Price         decimal.Decimal `json:"price"`
	Size          decimal.Decimal `json:"size"`
	FilledSize    decimal.Decimal `json:"filledSize"`
	RemainingSize decimal.Decimal `json:"remainingSize"`
	AvgFillPrice  decimal.Decimal `json:"avgFillPrice"`
	Status        Status          `json:"status"`
	CreatedAt     time.Time       `json:"createdAt"`
	ReduceOnly    bool            `json:"reduceOnly"`
	IOC           bool            `json:"ioc"`
	PostOnly      bool            `json:"postOnly"`
	Future        string          `json:"future"`
	ClientID      string          `json:"clientId"`
}

type OrdersHistoryParams struct {
	Market    *string `json:"market"`
	Limit     *int    `json:"limit"`
	StartTime *int    `json:"start_time"`
	EndTime   *int    `json:"end_time"`
}

type TriggerOrder struct {
	ID               int64            `json:"id"`
	OrderID          int64            `json:"orderId"`
	Market           string           `json:"market"`
	CreatedAt        time.Time        `json:"createdAt"`
	Error            string           `json:"error"`
	Future           string           `json:"future"`
	OrderPrice       decimal.Decimal  `json:"orderPrice"`
	ReduceOnly       bool             `json:"reduceOnly"`
	Side             Side             `json:"side"`
	Size             decimal.Decimal  `json:"size"`
	Status           Status           `json:"status"`
	TrailStart       decimal.Decimal  `json:"trailStart"`
	TrailValue       decimal.Decimal  `json:"trailValue"`
	TriggerPrice     decimal.Decimal  `json:"triggerPrice"`
	TriggeredAt      time.Time        `json:"triggeredAt"`
	Type             TriggerOrderType `json:"type"`
	OrderType        OrderType        `json:"orderType"`
	FilledSize       decimal.Decimal  `json:"filledSize"`
	AvgFillPrice     decimal.Decimal  `json:"avgFillPrice"`
	OrderStatus      string           `json:"orderStatus"`
	RetryUntilFilled bool             `json:"retryUntilFilled"`
}

type OpenTriggerOrdersParams struct {
	Market *string           `json:"market"`
	Type   *TriggerOrderType `json:"type"`
}

type Trigger struct {
	Error      string    `json:"error"`
	FilledSize float64   `json:"filledSize"`
	OrderSize  float64   `json:"orderSize"`
	OrderID    int64     `json:"orderId"`
	Time       time.Time `json:"time"`
}

type TriggerOrdersHistoryParams struct {
	Market    *string `json:"market"`
	Limit     *int    `json:"limit"`
	StartTime *int    `json:"start_time"`
	EndTime   *int    `json:"end_time"`
	Side      *string `json:"side"`
	Type      *string `json:"type"`
	OrderType *string `json:"orderType"`
}

type OrderParams struct {
	Market     *string          `json:"market"`
	Side       *string          `json:"side"`
	Price      *decimal.Decimal `json:"price"`
	Type       *string          `json:"type"`
	Size       *decimal.Decimal `json:"size"`
	ReduceOnly *bool            `json:"reduceOnly,omitempty"`
	IOC        *bool            `json:"ioc,omitempty"`
	PostOnly   *bool            `json:"postOnly,omitempty"`
	ClientID   *string          `json:"clienId,omitempty"`
}

type TriggerOrderParams struct {
	Market       *string          `json:"market"`
	Side         *string          `json:"side"`
	Size         *decimal.Decimal `json:"size"`
	Type         *string          `json:"type"`
	TriggerPrice *decimal.Decimal `json:"triggerPrice"`
	OrderPrice   *decimal.Decimal `json:"orderPrice"`
	ReduceOnly   *bool            `json:"reduceOnly,omitempty"`
	TrailValue   *decimal.Decimal `json:"trailValue,omitempty"`
}

type ModifyOrderParams struct {
	Price    *decimal.Decimal `json:"price,omitempty"`
	Size     *decimal.Decimal `json:"size,omitempty"`
	ClientID *string          `json:"clientId,omitempty"`
}

type ModifyTriggerOrderParams struct {
	TriggerPrice *decimal.Decimal `json:"triggerPrice,omitempty"`
	Size         *decimal.Decimal `json:"size,omitempty"`
	OrderPrice   *decimal.Decimal `json:"orderPrice,omitempty"`
	TrailValue   *decimal.Decimal `json:"trailValue,omitempty"`
}

type CancelAllParams struct {
	Market                *string `json:"market"`
	Side                  *string `json:"side"`
	ConditionalOrdersOnly *bool   `json:"conditionalOrdersOnly,omitempty"`
	LimitOrdersOnly       *bool   `json:"limitOrdersOnly"`
}
