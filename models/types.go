package models

import (
	"encoding/json"
	"fmt"
	"math"
	"time"
)

type Resolution int

const (
	Sec15    = 15
	Minute   = 60
	Minute5  = 300
	Minute15 = 900
	Hour     = 3600
	Hour4    = 14400
	Day      = 86400
)

type NumberTimeLimit struct {
	Limit     *int   `json:"limit,omitempty"`
	StartTime *int64 `json:"start_time,omitempty"`
	EndTime   *int64 `json:"end_time,omitempty"`
}

type Channel string

const (
	OrderBookChannel = Channel("orderbook")
	TradesChannel    = Channel("trades")
	TickerChannel    = Channel("ticker")
	MarketsChannel   = Channel("markets")
	FillsChannel     = Channel("fills")
	OrdersChannel    = Channel("orders")
)

type Operation string

const (
	Subscribe   = Operation("subscribe")
	UnSubscribe = Operation("unsubscribe")
)

type ResponseType string

const (
	Error        = ResponseType("error")
	Subscribed   = ResponseType("subscribed")
	UnSubscribed = ResponseType("unsubscribed")
	Info         = ResponseType("info")
	Partial      = ResponseType("partial")
	Update       = ResponseType("update")
)

type TransferStatus string

const Complete = TransferStatus("complete")

type OrderType string

const (
	LimitOrder  = OrderType("limit")
	MarketOrder = OrderType("market")
)

type Side string

const (
	Sell = Side("sell")
	Buy  = Side("buy")
)

type Status string

const (
	New    = Status("new")
	Open   = Status("open")
	Closed = Status("closed")
)

type TriggerOrderType string

const (
	Stop         = TriggerOrderType("stop")
	TrailingStop = TriggerOrderType("trailing_stop")
	TakeProfit   = TriggerOrderType("take_profit")
)

type DepositMethod string

const (
	Erc20 = DepositMethod("erc20")
	Trc20 = DepositMethod("trc20")
	Spl   = DepositMethod("spl")
	Omni  = DepositMethod("omni")
	Bep2  = DepositMethod("bep2")
)

type OptionType string

const (
	Call = OptionType("call")
	Put  = OptionType("put")
)

type LiquidityType string

const (
	Maker = LiquidityType("maker")
	Taker = LiquidityType("taker")
)

type UnstakeRequestStatus string

const (
	Pending   = UnstakeRequestStatus("pending")
	Cancelled = UnstakeRequestStatus("cancelled")
	Processed = UnstakeRequestStatus("processed")
)

type Result struct {
	Result string `json:"result"`
}

type FTXTime struct {
	Time time.Time
}

var ErrNilPtr = fmt.Errorf("nil pointer")

func (f *FTXTime) UnmarshalJSON(data []byte) error {
	var t float64
	if err := json.Unmarshal(data, &t); err != nil {
		return err
	}

	sec, nsec := math.Modf(t)
	f.Time = time.Unix(int64(sec), int64(nsec))
	return nil
}

func (f FTXTime) MarshalJSON() ([]byte, error) {
	return json.Marshal(float64(f.Time.UnixNano()) / float64(1000000000))
}
