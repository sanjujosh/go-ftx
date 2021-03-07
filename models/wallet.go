package models

import (
	"time"

	"github.com/shopspring/decimal"
)

type Coin struct {
	Bep2Asset     string          `json:"bep2Asset"`
	CanConvert    bool            `json:"canConvert"`
	CanDeposit    bool            `json:"canDeposit"`
	CanWithdraw   bool            `json:"canWithdraw"`
	Collateral    bool            `json:"collateral"`
	CreditTo      string          `json:"creditTo"`
	Erc20Contract string          `json:"erc20Contract"`
	Fiat          bool            `json:"fiat"`
	HasTag        bool            `json:"hasTag"`
	ID            string          `json:"id"`
	IsToken       bool            `json:"isToken"`
	Methods       []DepositMethod `json:"methods"`
	Name          string          `json:"name"`
	SplMint       string          `json:"splMint"`
	Trc20Contract string          `json:"trc20Contract"`
	UsdFungible   bool            `json:"usdFungible"`
}

type Balance struct {
	Coin  string          `json:"coin"`
	Free  decimal.Decimal `json:"free"`
	Total decimal.Decimal `json:"total"`
}

type DepositAddress struct {
	Address string `json:"address"`
	Tag     string `json:"tag"`
}

type CoinTxBase struct {
	Coin   string          `json:"coin"`
	ID     int64           `json:"id"`
	Size   decimal.Decimal `json:"size"`
	Status string          `json:"status"`
	Time   time.Time       `json:"time"`
}

type WalletTransaction struct {
	*CoinTxBase
	Fee   decimal.Decimal `json:"fee"`
	Notes string          `json:"notes"`
	Txid  string          `json:"txid"`
}

type DepositHistoryParams NumberTimeLimit

type Deposit struct {
	*WalletTransaction
	Confirmations int       `json:"confirmations"`
	ConfirmedTime time.Time `json:"confirmedTime"`
	SentTime      time.Time `json:"sentTime"`
}

type WithdrawalHistoryParams DepositHistoryParams

type Withdrawal struct {
	*WalletTransaction
	Address string `json:"address"`
	Tag     string `json:"tag"`
}

type RequestWithdrawalParams struct {
	Address  *string          `json:"address"`
	Code     *string          `json:"code,omitempty"`
	Coin     *string          `json:"coin"`
	Password *string          `json:"password,omitempty"`
	Size     *decimal.Decimal `json:"size"`
	Tag      *string          `json:"tag,omitempty"`
}

type AirDropParams NumberTimeLimit

type AirDrop struct {
	Coin   string          `json:"coin"`
	ID     int64           `json:"id"`
	Size   decimal.Decimal `json:"size"`
	Status string          `json:"status"`
	Time   time.Time       `json:"time"`
}

type SavedAddressParams struct {
	Address      *string `json:"address"`
	AddressName  *string `json:"addressName"`
	Coin         *string `json:"coin"`
	IsPrimetrust *bool   `json:"isPrimetrust"`
	Tag          *string `json:"tag,omitempty"`
}

type SavedAddress struct {
	Address          string    `json:"address"`
	Coin             string    `json:"coin"`
	Fiat             bool      `json:"fiat"`
	ID               int64     `json:"id"`
	IsPrimetrust     bool      `json:"isPrimetrust"`
	LastUsedAt       time.Time `json:"lastUsedAt"`
	Name             string    `json:"name"`
	Tag              string    `json:"tag"`
	Whitelisted      bool      `json:"whitelisted"`
	WhitelistedAfter string    `json:"whitelistedAfter"`
}
