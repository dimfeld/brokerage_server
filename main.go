package main

import (
	"fmt"
	"os"

	"github.com/inconshreveable/log15"

	"github.com/dimfeld/brokerage_server/brokers"
	"github.com/dimfeld/brokerage_server/server"
)

func main() {
	config, err := ReadConfig()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading config file: %s\n", err.Error())
		os.Exit(1)
	}

	// Right now only a single plugin is supported. In the future multiple plugins will be usable concurrently.
	ibConfig := config.Brokers["ib"]
	if ibConfig == nil {
		fmt.Fprintf(os.Stderr, "No config found at brokers.ib\n")
		os.Exit(1)
	}

	plugin, err := brokers.NewBrokerEngine(log15.Root(), "ib", ibConfig)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating plugin %s: %s\n", "ib", err.Error())
		os.Exit(1)
	}

	if err = plugin.Connect(); err != nil {
		fmt.Fprintf(os.Stderr, "Error connecting to broker %s: %s\n", "ib", err.Error())
		os.Exit(1)
	}

	server.Run(config.Bind, plugin)
}
