package types

import (
	"fmt"
	"time"
)

// VendorSpecific holds information that isn't common to the supported platforms,
// and isn't vital, but might be interesting to use when it's present.
type VendorSpecific struct {
	Data map[string]string
	Keys []string
}

type Account struct {
	Id          string
	Broker      string
	Description string
}

func (a Account) String() string {
	return fmt.Sprintf("%s:%s - %s", a.Broker, a.Id, a.Description)
}

type ConnectionStatus struct {
	Connected bool
	Error     error
}

type Tristate int

const (
	Yes Tristate = iota
	No
	Maybe
)

type Quote struct {
	High   float64
	Low    float64
	Open   float64
	Close  float64
	Volume float64
	Mark   float64

	Bid     float64
	BidSize float64
	BidExch string

	Ask     float64
	AskSize float64
	AskExch string

	LastTime time.Time
	Last     float64
	LastSize float64
	LastExch string

	OptionOpenInterest  float64
	OptionHistoricalVol float64
	OptionImpliedVol    float64
	OptionCallOpenInt   float64
	OptionPutOpenInt    float64
	OptionCallVolume    float64
	OptionPutVolume     float64

	Shortable Tristate

	AvgVol float64 // Not supported by all brokers

	YearHigh float64
	YearLow  float64

	Time time.Time
}

type OptionQuote struct {
	Quote
	Delta float64
	Gamma float64
	Theta float64
	Vega  float64
	Rho   float64 // Not always supported
}

type SymbolType int

const (
	SymbolEquity SymbolType = iota
	SymbolOption
)

// func (t SymbolType) String() string {
// 	switch t {
// 	case SymbolEquity:
// 		return "Equity"
// 	case SymbolOption:
// 		return "Option"
// 	default:
// 		return "Unknown"
// 	}
// }

type SymbolDetails struct {
	Symbol      string
	Description string
	Vendor      VendorSpecific
}
