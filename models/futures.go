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
	VolumeUsd24h        float64         `json:"volumeUsd24h"`
	Volume              float64         `json:"volume"`
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
	SizeIncrement       decimal.Decimal `json:"sizeIncrement"`
	Underlying          string          `json:"underlying"`
	UpperBound          decimal.Decimal `json:"upperBound"`
	Type                string          `json:"type"`
}

type FutureStats struct {
	Volume                   decimal.Decimal `json:"volume"`
	NextFundingRate          float64         `json:"nextFundingRate"`
	NextFundingTime          time.Time       `json:"nextFundingTime"`
	ExpirationPrice          decimal.Decimal `json:"expirationPrice"`
	PredictedExpirationPrice decimal.Decimal `json:"predictedExpirationPrice"`
	StrikePrice              decimal.Decimal `json:"strikePrice"`
	OpenInterest             float64         `json:"openInterest"`
}

type FundingRatesParams struct {
	StartTime *int64 `json:"start_time,omitempty"`
	EndTime   *int64 `json:"end_time,omitempty"`
}

type FundingRates struct {
	Future string    `json:"future"`
	Rate   float64   `json:"rate"`
	Time   time.Time `json:"time"`
}

type FutureExpired struct {
	Ask                   decimal.Decimal `json:"ask"`
	Bid                   decimal.Decimal `json:"bid"`
	Description           string          `json:"description"`
	Enabled               bool            `json:"enabled"`
	Expired               bool            `json:"expired"`
	Expiry                time.Time       `json:"expiry"`
	ExpiryDescription     string          `json:"expiryDescription"`
	Group                 string          `json:"group"`
	ImfFactor             float64         `json:"imfFactor"`
	Index                 float64         `json:"index"`
	Last                  decimal.Decimal `json:"last"`
	LowerBound            decimal.Decimal `json:"lowerBound"`
	MarginPrice           float64         `json:"marginPrice"`
	Mark                  decimal.Decimal `json:"mark"`
	MoveStart             string          `json:"moveStart"`
	Name                  string          `json:"name"`
	Perpetual             bool            `json:"perpetual"`
	PositionLimitWeight   float64         `json:"positionLimitWeight"`
	PostOnly              bool            `json:"postOnly"`
	PriceIncrement        decimal.Decimal `json:"priceIncrement"`
	SizeIncrement         decimal.Decimal `json:"sizeIncrement"`
	Type                  string          `json:"type"`
	Underlying            string          `json:"underlying"`
	UnderlyingDescription string          `json:"underlyingDescription"`
	UpperBound            decimal.Decimal `json:"upperBound"`
}

type HistoricalIndexParams struct {
	Resolution *int `json:"resolution"`
	Limit      *int `json:"limit,omitempty"`
	StartTime  *int `json:"start_time,omitempty"`
	EndTime    *int `json:"end_time,omitempty"`
}

type HistoricalIndex struct {
	StartTime time.Time       `json:"startTime"`
	Open      decimal.Decimal `json:"open"`
	High      decimal.Decimal `json:"high"`
	Low       decimal.Decimal `json:"low"`
	Close     decimal.Decimal `json:"close"`
	Volume    float64         `json:"volume"`
}
