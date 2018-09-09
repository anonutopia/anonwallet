package main

import (
	"bytes"
	"encoding/json"
	"log"
	"os"
	"os/exec"
	"strconv"
)

type BitcoinGenerator struct {
}

func (bg *BitcoinGenerator) getAddress() (string, error) {
	cmd := exec.Command("/usr/local/bin/electrum", "createnewaddress")
	cmd.Env = append(os.Environ(), "HOME=/home/kriptokuna")
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		log.Println("Error in BitcoinGenerator.getAddress: %s" + string(stderr.Bytes()))
		return "", err
	}
	return string(stdout.Bytes()), nil
}

func (bg *BitcoinGenerator) getBalance(address string) (float64, error) {
	cmd := exec.Command("/usr/local/bin/electrum", "getaddressbalance", address)
	cmd.Env = append(os.Environ(), "HOME=/home/kriptokuna")
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		log.Println("Error in BitcoinGenerator.getBalance: %s" + string(stderr.Bytes()))
		return 0, err
	}

	br := &BalanceResponse{}
	err = json.Unmarshal(stdout.Bytes(), br)
	if err != nil {
		return 0, err
	}

	balance, err := strconv.ParseFloat(br.Confirmed, 64)
	if err != nil {
		return 0, err
	}

	return balance, nil
}

func initBtcGen() *BitcoinGenerator {
	bg := &BitcoinGenerator{}
	return bg
}

type BalanceResponse struct {
	Confirmed   string `json:"confirmed"`
	Unconfirmed string `json:"unconfirmed"`
}
