package types

import (
	"context"
	"errors"
)

var ERR_PLUGIN_NOT_IMPLEMENTED = errors.New("Not implemented for this broker")

type DebugLevel int

const (
	DEBUG_OFF     DebugLevel = 0
	DEBUG_NORMAL             = 1
	DEBUG_VERBOSE            = 2
	DEBUG_TRACE              = 3
)

type BrokerageServerPluginV1 interface {
	Connect() error
	Close() error
	Status() *ConnectionStatus
	Error() error
	SetDebugLevel(level DebugLevel)

	// Accounts(ctx context.Context) ([]*Account, error)
	GetStockQuote(ctx context.Context, symbol string) (*Quote, error)
	GetOptionsChain(ctx context.Context, symbol string) (OptionChain, error)
}
