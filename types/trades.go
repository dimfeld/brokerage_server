package types

import "time"

type Trade struct {
	Account string `json:"account"`
	Broker  string `json:"broker"`
	TradeId string `json:"trade_id"`

	Symbol     string  `json:"symbol"`
	OptionType string  `json:"type,omitempty"`
	Strike     float64 `json:"strike,omitempty"`
	Expiration string  `json:"expiration,omitempty"`
	Multiplier string  `json:"multiplier,omitempty"`

	Exchange    string  `json:"exchange"`
	Side        string  `json:"side"`
	Size        int     `json:"size"`
	Price       float64 `json:"price"`
	AvgPrice    float64 `json:"avg_price"`
	CumQty      int64   `json:"cumulative_size"`
	Commissions float64 `json:"commissions"`
	RealizedPnL float64 `json:"realized_pnl,omitempty"`

	Time    time.Time   `json:"time"`
	RawData interface{} `json:"raw_data,omitempty"`
}
