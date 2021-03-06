package types

import (
	"context"
	"errors"
	"time"
)

var ERR_PLUGIN_NOT_IMPLEMENTED = errors.New("Not implemented for this broker")

type DebugLevel int

const (
	DEBUG_OFF     DebugLevel = 0
	DEBUG_NORMAL             = 1
	DEBUG_VERBOSE            = 2
	DEBUG_TRACE              = 3
)

type OptionsQuoteParams struct {
	Underlying  string
	Expirations []string
	Strikes     []float64
	Puts        bool
	Calls       bool
}

type BrokerageServerPluginV1 interface {
	Connect() error
	Close() error
	Status() *ConnectionStatus
	Error() error
	SetDebugLevel(level DebugLevel)

	AccountList(ctx context.Context) ([]*Account, error)
	GetStockQuote(ctx context.Context, symbol string) (*Quote, error)
	GetOptionsChain(ctx context.Context, symbol string) (OptionChain, error)
	GetOptionsQuotes(ctx context.Context, params OptionsQuoteParams) ([]*OptionQuote, error)
	GetHistoricalData(ctx context.Context, params HistoricalDataParams) ([]*Quote, error)
	GetPositions(ctx context.Context) ([]*Position, error)
	GetTrades(ctx context.Context, startDate time.Time) ([]*Trade, error)
}

type BrokerageServerPluginLatest = BrokerageServerPluginV1
