package brokers

import (
	"fmt"

	ib "github.com/dimfeld/brokerage_server_ib"
	plugin_intf "github.com/dimfeld/brokerage_server_plugin_intf"
)

type BrokerEngine struct {
	plugin_intf.BrokerageServerPluginV1
	Name string
}

func getPlugin(name string) (plugin_intf.BrokerageServerPluginV1, error) {
	var plugin interface{}
	switch name {
	case "ib":
		plugin = ib.New()
	default:
		return nil, fmt.Errorf("Unwknown plugin %s", name)
	}

	return wrapPlugin(name, plugin)
}

// Retunr
func wrapPlugin(name string, p interface{}) (plugin_intf.BrokerageServerPluginV1, error) {
	switch p.(type) {
	case plugin_intf.BrokerageServerPluginV1:
		return p.(plugin_intf.BrokerageServerPluginV1), nil
	default:
		return nil, fmt.Errorf("Plugin %s does not conform to the plugin interface", name)
	}
}

func NewBrokerEngine(name string) (*BrokerEngine, error) {
	var err error
	engine := &BrokerEngine{
		Name: name,
	}

	if engine.BrokerageServerPluginV1, err = getPlugin(name); err != nil {
		return nil, err
	}

	return engine, nil
}
