package types

type Account struct {
	ID     string `json:"id"`
	Broker string `json:"broker"`

	Type           string `json:"type"`
	NetLiquidation string `json:"net_liquidation"`
	BuyingPower    string `json:"buying_power"`
	MarginReq      string `json:"margin_req"`
	AvailableFunds string `json:"available_funds"`
}
