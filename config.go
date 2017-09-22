package main

import (
	"encoding/json"
	"os"

	jsoniter "github.com/json-iterator/go"
)

type ServerConfig struct {
	Brokers map[string]json.RawMessage `json:"brokers"`
}

func NewConfig() (*ServerConfig, error) {
	file, err := os.Open("config.json")
	if err != nil {
		return nil, err
	}

	defer file.Close()

	config := &ServerConfig{}
	err = jsoniter.NewDecoder(file).Decode(config)
	return config, err
}
