package types

type Account struct {
	ID     string `json:"id"`
	Broker string `json:"broker"`

	Type           string `json:"type"`
	NetLiquidation string `json:"net_liquidation"`
	BuyingPower    string `json:"buying_power"`
	MarginReq      string `json:"margin_req"`
	AvailableFunds string `json:"available_funds"`
	RealizedPnL    string `json:"realized_pnl"`
	UnrealizedPnL  string `json:"unrealized_pnl"`
}

type Position struct {
	Symbol     string    `json:"symbol"`
	Strike     float64   `json:"strike,omitempty"`
	Multiplier string    `json:"multiplier,omitempty"`
	Expiration string    `json:"expiration,omitempty"`
	PutOrCall  PutOrCall `json:"option_type,omitempty"`

	Position      int     `json:"position"`
	MarketPrice   float64 `json:"mkt_price"`
	MarketValue   float64 `json:"mkt_value"`
	AvgPrice      float64 `json:"avg_price"`
	UnrealizedPnL float64 `json:"unrealized_pnl"`
	RealizedPnL   float64 `json:"realized_pnl"`

	RawData interface{} `json:"raw,omitempty"`

	Broker  string `json:"broker"`
	Account string `json:"account"`
}
