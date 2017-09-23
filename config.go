package main

import (
	"encoding/json"
	"os"

	log "github.com/inconshreveable/log15"
	jsoniter "github.com/json-iterator/go"
)

type ServerConfig struct {
	Brokers    map[string]json.RawMessage `json:"brokers"`
	Bind       string                     `json:"bind"`
	Production bool                       `json:"production"`
	Debug      bool                       `json:"debug"`
}

func configureLogging(config *ServerConfig) {
	var formatter log.Format
	if config.Production {
		formatter = log.JsonFormat()
	} else {
		formatter = log.TerminalFormat()
	}

	logLevel := log.LvlInfo
	if config.Debug {
		logLevel = log.LvlDebug
	}
	log.Root().SetHandler(log.LvlFilterHandler(logLevel, log.StreamHandler(os.Stdout, formatter)))
}

func ReadConfig() (*ServerConfig, error) {
	file, err := os.Open("config.json")
	if err != nil {
		return nil, err
	}

	defer file.Close()

	config := &ServerConfig{}
	err = jsoniter.NewDecoder(file).Decode(config)

	if len(config.Bind) == 0 {
		config.Bind = "127.0.0.1:6543"
	}

	configureLogging(config)
	return config, err
}
