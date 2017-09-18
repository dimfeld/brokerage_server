package types

import "errors"

var ERR_PLUGIN_NOT_IMPLEMENTED = errors.New("Not implemented for this broker")

type BrokerageServerPluginV1 interface {
	Connect() error
	Close() error
	Status() *ConnectionStatus
	Error() error

	Accounts() []*Account
	GetStockQuote(symbol string) *Quote
}
