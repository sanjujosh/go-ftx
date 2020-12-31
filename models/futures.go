package models

import (
	"time"

	"github.com/shopspring/decimal"
)

type Future struct {
	Ask                 decimal.Decimal `json:"ask"`
	Bid                 decimal.Decimal `json:"bid"`
	Change1h            decimal.Decimal `json:"change1h"`
	Change24h           decimal.Decimal `json:"change24h"`
	ChangeBod           decimal.Decimal `json:"changeBod"`
	VolumeUsd24h        decimal.Decimal `json:"volumeUsd24h"`
	Volume              decimal.Decimal `json:"volume"`
	Description         string          `json:"description"`
	Enabled             bool            `json:"enabled"`
	Expired             bool            `json:"expired"`
	Expiry              time.Time       `json:"expiry"`
	Index               float64         `json:"index"`
	ImfFactor           float64         `json:"imfFactor"`
	Last                decimal.Decimal `json:"last"`
	LowerBound          decimal.Decimal `json:"lowerBound"`
	Mark                decimal.Decimal `json:"mark"`
	Name                string          `json:"name"`
	Perpetual           bool            `json:"perpetual"`
	PositionLimitWeight float64         `json:"positionLimitWeight"`
	PostOnly            bool            `json:"postOnly"`
	PriceIncrement      decimal.Decimal `json:"priceIncrement"`
	SizeIncrement       float64         `json:"sizeIncrement"`
	Underlying          string          `json:"underlying"`
	UpperBound          decimal.Decimal `json:"upperBound"`
	Type                string          `json:"type"`
}
