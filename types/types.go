package types

import (
	"errors"
	"net/http"
)

// VendorSpecific holds information that isn't common to the supported platforms,
// and isn't vital, but might be interesting to use when it's present.
type VendorSpecific struct {
	Data map[string]string `json:"data"`
	// Keys defines a preferred order to print the keys.
	Keys []string `json:"keys"`
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

type PutOrCall string

const (
	Put  PutOrCall = "PUT"
	Call           = "CALL"
)

type OptionChain struct {
	Underlying  string    `json:"underlying"`
	Multiplier  string    `json:"multiplier,omitempty"`
	Exchanges   []string  `json:"exchanges,omitempty"`
	Strikes     []float64 `json:"strikes"`
	Expirations []string  `json:"expirations"`
}

type Option struct {
	Underlying string
	Strike     float64
	Expiration string
	Type       PutOrCall
}

type OptionCombo struct {
	Legs []Option
}

type ErrorWithCode struct {
	error
	code int
}

func (ec ErrorWithCode) Code() int {
	return ec.code
}

var (
	ErrSymbolNotFound = ErrorWithCode{errors.New("symbol not found"), http.StatusNotFound}
	ErrDisconnected   = errors.New("broker disconnected")
)

func ArgError(message string) ErrorWithCode {
	return ErrorWithCode{
		error: errors.New(message),
		code:  http.StatusBadRequest,
	}
}
