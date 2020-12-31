package models

import (
	"time"

	"github.com/shopspring/decimal"
)

type Future struct {
	Ask                 decimal.Decimal `json:"ask,omitempty"`
	Bid                 decimal.Decimal `json:"bid,omitempty"`
	Change1h            decimal.Decimal `json:"change1h,omitempty"`
	Change24h           decimal.Decimal `json:"change24h,omitempty"`
	ChangeBod           decimal.Decimal `json:"changeBod,omitempty"`
	VolumeUsd24h        float64         `json:"volumeUsd24h,omitempty"`
	Volume              float64         `json:"volume,omitempty"`
	Description         string          `json:"description,omitempty"`
	Enabled             bool            `json:"enabled,omitempty"`
	Expired             bool            `json:"expired,omitempty"`
	Expiry              time.Time       `json:"expiry,omitempty"`
	Index               float64         `json:"index,omitempty"`
	ImfFactor           float64         `json:"imfFactor,omitempty"`
	Last                decimal.Decimal `json:"last,omitempty"`
	LowerBound          decimal.Decimal `json:"lowerBound,omitempty"`
	Mark                decimal.Decimal `json:"mark,omitempty"`
	Name                string          `json:"name,omitempty"`
	Perpetual           bool            `json:"perpetual,omitempty"`
	PositionLimitWeight float64         `json:"positionLimitWeight,omitempty"`
	PostOnly            bool            `json:"postOnly,omitempty"`
	PriceIncrement      decimal.Decimal `json:"priceIncrement,omitempty"`
	SizeIncrement       float64         `json:"sizeIncrement,omitempty"`
	Underlying          string          `json:"underlying,omitempty"`
	UpperBound          decimal.Decimal `json:"upperBound,omitempty"`
	Type                string          `json:"type,omitempty"`
}
