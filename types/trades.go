package types

import "time"

type Execution struct {
	ExecutionId string  `json:"id"`
	OptionType  string  `json:"type,omitempty"`
	Strike      float64 `json:"strike,omitempty"`
	Expiration  string  `json:"expiration,omitempty"`
	Multiplier  string  `json:"multiplier,omitempty"`

	Exchange    string  `json:"exchange"`
	Side        string  `json:"side"`
	Size        int     `json:"size"`
	Price       float64 `json:"price"`
	AvgPrice    float64 `json:"avg_price"`
	Commissions float64 `json:"commissions"`
	RealizedPnL float64 `json:"realized_pnl,omitempty"`

	Time    time.Time   `json:"time"`
	RawData interface{} `json:"raw_data,omitempty"`
}

type Trade struct {
	Account string `json:"account"`
	Broker  string `json:"broker"`
	OrderId string `json:"id"`
	Symbol  string `json:"symbol"`

	Executions []*Execution `json:"executions"`

	Time    time.Time   `json:"time"`
	RawData interface{} `json:"raw_data,omitempty"`
}
