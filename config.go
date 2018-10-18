package main

import (
	"encoding/json"
	"log"
	"os"
)

type config struct {
	EthNetwork      uint   `json:"eth_network"`
	NodeAddress     string `json:"node_address"`
	Port            uint   `json:"port"`
	Debug           bool   `json:"debug"`
	BtcMasterKey    string `json:"btc_master_key"`
	DbName          string `json:"db_name"`
	DbUser          string `json:"db_user"`
	DbPass          string `json:"db_pass"`
	SendgridKey     string `json:"sendgrid_key"`
	WavesNodeApiKey string `json:"wavesnode_apikey"`
}

func (sc *config) Load(configFile string) error {
	file, err := os.Open(configFile)

	if err != nil {
		log.Printf("[Config.Load] Got error while opening config file: %v", err)
		return err
	}

	decoder := json.NewDecoder(file)

	err = decoder.Decode(&sc)

	if err != nil {
		log.Printf("[Config.Load] Got error while decoding JSON: %v", err)
		return err
	}

	return nil
}

func initConfig() *config {
	c := &config{}
	c.Load("config.json")
	return c
}
