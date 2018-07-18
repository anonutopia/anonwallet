package main

import (
	"fmt"
	"log"

	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcutil/hdkeychain"
)

type BitcoinGenerator struct {
	masterKey *hdkeychain.ExtendedKey
}

func (bg *BitcoinGenerator) getAddress(userID uint32) (string, error) {
	acc, err := bg.masterKey.Child(hdkeychain.HardenedKeyStart + userID)
	if err != nil {
		return "", err
	}

	accExt, err := acc.Child(0)
	if err != nil {
		return "", err
	}

	accExt0, err := accExt.Child(0)
	if err != nil {
		return "", err
	}

	accExt0Addr, err := accExt0.Address(&chaincfg.MainNetParams)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s", accExt0Addr), nil
}

func initBtcGen() *BitcoinGenerator {
	var err error
	bg := &BitcoinGenerator{}

	bg.masterKey, err = hdkeychain.NewKeyFromString(conf.BtcMasterKey)
	if err != nil {
		log.Fatalf("hdkeychain.NewKeyFromString error: %s", err)
	}

	return bg
}
