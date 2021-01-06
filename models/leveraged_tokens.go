package models

import (
	"time"

	"github.com/shopspring/decimal"
)

type LeveragedToken struct {
	Name             string          `json:"name"`
	Description      string          `json:"description"`
	Underlying       string          `json:"underlying"`
	Leverage         float64         `json:"leverage"`
	Outstanding      decimal.Decimal `json:"outstanding"`
	PricePerShare    decimal.Decimal `json:"pricePerShare"`
	PositionPerShare decimal.Decimal `json:"positionPerShare"`
	UnderlyingMark   decimal.Decimal `json:"underlyingMark"`
	ContractAddress  string          `json:"contractAddress"`
	Change1h         decimal.Decimal `json:"change1h"`
	Change24h        decimal.Decimal `json:"change24h"`
}

type TokenInfo LeveragedToken

type LeveragedTokenBalance struct {
	Token   string          `json:"token"`
	Balance decimal.Decimal `json:"balance"`
}

type LeveragedTokenCreationRequest struct {
	ID            int64           `json:"id"`
	Token         string          `json:"token"`
	RequestedSize decimal.Decimal `json:"requestedSize"`
	Pending       bool            `json:"pending"`
	CreatedSize   decimal.Decimal `json:"createdSize"`
	Price         decimal.Decimal `json:"price"`
	Cost          decimal.Decimal `json:"cost"`
	Fee           decimal.Decimal `json:"fee"`
	RequestedAt   time.Time       `json:"requestedAt"`
	FulfilledAt   time.Time       `json:"fulfilledAt"`
}

type LeveragedTokenCreation struct {
	ID            int64           `json:"id"`
	Token         string          `json:"token"`
	RequestedSize decimal.Decimal `json:"requestedSize"`
	Cost          decimal.Decimal `json:"cost"`
	Pending       bool            `json:"pending"`
	RequestedAt   time.Time       `json:"requestedAt"`
}

type LeveragedTokenRedemptionRequest struct {
	ID          int64           `json:"id"`
	Token       string          `json:"token"`
	Size        decimal.Decimal `json:"size"`
	Pending     bool            `json:"pending"`
	Price       decimal.Decimal `json:"price"`
	Proceeds    decimal.Decimal `json:"proceeds"`
	Fee         decimal.Decimal `json:"fee"`
	RequestedAt time.Time       `json:"requestedAt"`
	FulfilledAt time.Time       `json:"fulfilledAt"`
}

type LeveragedTokenRedemption struct {
	ID                int64           `json:"id"`
	Token             string          `json:"token"`
	Size              decimal.Decimal `json:"size"`
	ProjectedProceeds decimal.Decimal `json:"projectedProceeds"`
	Pending           bool            `json:"pending"`
	RequestedAt       time.Time       `json:"requestedAt"`
}
