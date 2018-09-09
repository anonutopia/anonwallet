package main

import (
	"os/exec"
)

type BitcoinGenerator struct {
}

func (bg *BitcoinGenerator) getAddress() (string, error) {
	out, err := exec.Command("electrum createnewaddress").Output()
	if err != nil {
		return "", err
	}
	return string(out), nil
}

func initBtcGen() *BitcoinGenerator {
	bg := &BitcoinGenerator{}
	return bg
}
