package brokers

import (
	"encoding/json"
	"fmt"

	"github.com/inconshreveable/log15"

	"github.com/dimfeld/brokerage_server/types"
	ib "github.com/dimfeld/brokerage_server_ib"
)

type BrokerEngine struct {
	types.BrokerageServerPluginV1
	Name string
}

func getPlugin(logger log15.Logger, name string, configData json.RawMessage) (types.BrokerageServerPluginV1, error) {
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
func wrapPlugin(name string, p interface{}) (types.BrokerageServerPluginV1, error) {
	switch p.(type) {
	case types.BrokerageServerPluginV1:
		return p.(types.BrokerageServerPluginV1), nil
	default:
		return nil, fmt.Errorf("Plugin %s does not conform to the plugin interface", name)
	}
}

func NewBrokerEngine(logger log15.Logger, name string, config json.RawMessage) (*BrokerEngine, error) {
	var err error
	engine := &BrokerEngine{
		Name: name,
	}

	if engine.BrokerageServerPluginV1, err = getPlugin(logger, name, config); err != nil {
		return nil, err
	}

	return engine, nil
}
