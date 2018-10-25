package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"os/exec"
	"strconv"
	"strings"
)

type BitcoinGenerator struct {
}

func (bg *BitcoinGenerator) getAddress() (string, error) {
	cmdStr := fmt.Sprintf("/usr/local/bin/electrum createnewaddress")
	cmd := exec.Command("bash", "-c", cmdStr)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()

	if err != nil {
		log.Println("Error in BitcoinGenerator.getAddress: " + string(stderr.Bytes()))
		return "", err
	}
	return strings.TrimRight(string(stdout.Bytes()), "\n"), nil
}

func (bg *BitcoinGenerator) getBalance(address string) (float64, error) {
	cmdStr := fmt.Sprintf("/usr/local/bin/electrum getaddressbalance %s", address)
	cmd := exec.Command("bash", "-c", cmdStr)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()

	if err != nil {
		log.Printf("Error in BitcoinGenerator.getBalance (address: %s): %s", address, string(stderr.Bytes()))
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
