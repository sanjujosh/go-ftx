package models

import (
	"time"

	"github.com/shopspring/decimal"
)

type Option struct {
	Underlying string          `json:"underlying"`
	Type       OptionType      `json:"type"`
	Strike     decimal.Decimal `json:"strike"`
	Expiry     time.Time       `json:"expiry"`
}

type OptionOutput Option

type OptionInput struct {
	Underlying *string          `json:"underlying"`
	Type       *OptionType      `json:"type"`
	Strike     *decimal.Decimal `json:"strike"`
	Expiry     *int64           `json:"expiry"`
}

type OptionQuoteRequest struct {
	ID             int64           `json:"id"`
	Option         Option          `json:"option"`
	Side           Side            `json:"side"`
	Size           decimal.Decimal `json:"size"`
	Time           time.Time       `json:"time"`
	RequestExpiry  time.Time       `json:"requestExpiry"`
	Status         OrderStatus     `json:"status"`
	HideLimitPrice bool            `json:"hideLimitPrice"`
	LimitPrice     decimal.Decimal `json:"limitPrice"`
	Quotes         []OptionQuote   `json:"quotes"`
}

type OptionQuote struct {
	Collateral  decimal.Decimal `json:"collateral"`
	ID          int64           `json:"id"`
	Price       decimal.Decimal `json:"price"`
	QuoteExpiry time.Time       `json:"quoteExpiry"`
	Status      OrderStatus     `json:"status"`
	Time        time.Time       `json:"time"`
}

type OptionQuoteRequestParams struct {
	*OptionInput
	Side           *Side            `json:"side"`
	Size           *decimal.Decimal `json:"size"`
	LimitPrice     *decimal.Decimal `json:"limitPrice,omitempty"`
	HideLimitPrice *bool            `json:"hideLimitPrice,omitempty"`
	RequestExpiry  *int64           `json:"requestExpiry,omitempty"`
	CounterpartyID *int64           `json:"counterpartyId,omitempty"`
}

type CreateQuoteRequest struct {
	OptionInput
	Side Side            `json:"side"`
	Size decimal.Decimal `json:"size"`
}

type CancelQuoteRequest OptionQuoteRequest

type QuotesForOptionQuoteRequest struct {
	Collateral  decimal.Decimal `json:"collateral"`
	ID          int64           `json:"id"`
	Option      Option          `json:"option"`
	Price       decimal.Decimal `json:"price"`
	QuoteExpiry time.Time       `json:"quoteExpiry"`
	QuoterSide  Side            `json:"quoterSide"`
	RequestSide Side            `json:"requestSide"`
	Size        decimal.Decimal `json:"size"`
}

type UserOptionQuote struct {
	QuotesForOptionQuoteRequest
	Status OrderStatus `json:"status"`
	Time   time.Time   `json:"time"`
}

type AccountOptionsInfo struct {
	UsdBalance                   decimal.Decimal `json:"usdBalance"`
	LiquidationPrice             decimal.Decimal `json:"liquidationPrice"`
	Liquidated                   bool            `json:"liquidated"`
	MaintenanceMarginRequirement decimal.Decimal `json:"maintenanceMarginRequirement"`
	InitialMarginRequirement     decimal.Decimal `json:"initialMarginRequirement"`
}

type OptionPosition struct {
	EntryPrice            decimal.Decimal `json:"entryPrice"`
	NetSize               decimal.Decimal `json:"netSize"`
	Option                Option          `json:"option"`
	Side                  Side            `json:"side"`
	Size                  decimal.Decimal `json:"size"`
	PessimisticValuation  decimal.Decimal `json:"pessimisticValuation"`
	PessimisticIndexPrice decimal.Decimal `json:"pessimisticIndexPrice"`
	PessimisticVol        decimal.Decimal `json:"pessimisticVol"`
}

type PublicOptionTrade struct {
	ID     int64           `json:"id"`
	Size   decimal.Decimal `json:"size"`
	Price  decimal.Decimal `json:"price"`
	Option Option          `json:"option"`
	Time   time.Time       `json:"time"`
}

type OptionFill struct {
	Fee       decimal.Decimal `json:"fee"`
	FeeRate   decimal.Decimal `json:"feeRate"`
	ID        int64           `json:"id"`
	Liquidity LiquidityType   `json:"liquidity"`
	Option    Option          `json:"option"`
	Price     decimal.Decimal `json:"price"`
	QuoteID   int             `json:"quoteId"`
	Side      Side            `json:"side"`
	Size      decimal.Decimal `json:"size"`
	Time      time.Time       `json:"time"`
}

type OptionsVolume struct {
	Contracts       decimal.Decimal `json:"contracts"`
	UnderlyingTotal decimal.Decimal `json:"underlying_total"`
}

type OptionsHistoricalVolumes struct {
	NumContracts decimal.Decimal `json:"numContracts"`
	EndTime      time.Time       `json:"endTime"`
	StartTime    time.Time       `json:"startTime"`
}

type OptionsHistoricalOpenInterest struct {
	NumContracts decimal.Decimal `json:"numContracts"`
	Time         time.Time       `json:"time"`
}
