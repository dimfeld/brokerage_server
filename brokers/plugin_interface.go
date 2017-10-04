package brokers

import (
	"encoding/json"
	"fmt"

	"github.com/inconshreveable/log15"

	"github.com/dimfeld/brokerage_server/types"
	ib "github.com/dimfeld/brokerage_server_ib"
)

type BrokerEngine struct {
	types.BrokerageServerPluginLatest
	Name string
}

func getPlugin(logger log15.Logger, name string, configData json.RawMessage) (types.BrokerageServerPluginLatest, error) {
	var plugin interface{}
	var err error
	switch name {
	case "ib":
		plugin, err = ib.New(logger, configData)
	default:
		return nil, fmt.Errorf("Unwknown plugin %s", name)
	}

	if err != nil {
		return nil, err
	}

	return wrapPlugin(name, plugin)
}

// Return a plugin wrapped with dummy function to conform to the latest interface, where needed.
func wrapPlugin(name string, p interface{}) (types.BrokerageServerPluginLatest, error) {
	switch p.(type) {
	case types.BrokerageServerPluginV1:
		return p.(types.BrokerageServerPluginLatest), nil
	default:
		return nil, fmt.Errorf("Plugin %s does not conform to the plugin interface", name)
	}
}

func NewBrokerEngine(logger log15.Logger, name string, config json.RawMessage) (*BrokerEngine, error) {
	var err error
	engine := &BrokerEngine{
		Name: name,
	}

	if engine.BrokerageServerPluginLatest, err = getPlugin(logger, name, config); err != nil {
		return nil, err
	}

	return engine, nil
}

type Priorities struct {
	EquityQuotes []string `json:"equity_quotes"`
	OptionQuotes []string `json:"option_quotes"`
	EquityInfo   []string `json:"equity_info"`
	OptionInfo   []string `json:"option_info"`
	Other        []string `json:"other"`
	Fallback     string
}

type PriorityAndPurpose struct {
	Priority []string
	Purpose  string
}

type Purpose int

const (
	PurposeOther        Purpose = 0
	PurposeEquityQuotes         = 1
	PurposeEquityInfo           = 2
	PurposeOptionQuotes         = 3
	PurposeOptionInfo           = 4
)

func (p Priorities) ToList() []PriorityAndPurpose {
	expectedOrder := []Purpose{PurposeOther, PurposeEquityQuotes, PurposeEquityInfo, PurposeOptionQuotes, PurposeOptionInfo}
	for i, value := range expectedOrder {
		if i != int(value) {
			panic(fmt.Sprintf("Purpose list order mismatch, saw %d expected %d", value, i))
		}
	}

	return []PriorityAndPurpose{
		PriorityAndPurpose{Purpose: "other", Priority: p.Other},
		PriorityAndPurpose{Purpose: "equity_quotes", Priority: p.EquityQuotes},
		PriorityAndPurpose{Purpose: "equity_info", Priority: p.EquityInfo},
		PriorityAndPurpose{Purpose: "option_quotes", Priority: p.OptionQuotes},
		PriorityAndPurpose{Purpose: "option_info", Priority: p.OptionInfo},
	}
}

type EngineList struct {
	Engines    map[string]*BrokerEngine
	Priorities []PriorityAndPurpose
	Fallback   string
}

func NewEngineList(engines map[string]*BrokerEngine, priorities *Priorities) *EngineList {
	return &EngineList{
		Engines:    engines,
		Priorities: priorities.ToList(),
		Fallback:   priorities.Fallback,
	}
}

func (b *EngineList) Get(purpose Purpose) (*BrokerEngine, error) {

	if int(purpose) > len(b.Priorities) {
		return nil, fmt.Errorf("Unknown purpose %d", purpose)
	}
	prioList := b.Priorities[purpose]

	var engine *BrokerEngine
	for _, name := range prioList.Priority {
		if e, ok := b.Engines[name]; ok {
			engine = e
			break
		}
	}

	if engine == nil && len(b.Fallback) != 0 {
		engine = b.Engines[b.Fallback]
	}

	if engine == nil {
		return nil, fmt.Errorf("No suitable engine for %s operations. Please check your config.", prioList.Purpose)
	}

	return engine, nil
}
