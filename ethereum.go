package main

import (
	"github.com/ethereum/go-ethereum/accounts/keystore"
)

type EthereumGenerator struct {
	keystore *keystore.KeyStore
}

func (eg *EthereumGenerator) getAddress(userID uint32) (string, error) {
	account, err := eg.keystore.NewAccount(conf.BtcMasterKey)
	if err != nil {
		return "", err
	}

	return account.Address.Hex(), nil
}

func initEthGen() *EthereumGenerator {
	eg := &EthereumGenerator{}
	eg.keystore = keystore.NewKeyStore("./wallets", keystore.StandardScryptN, keystore.StandardScryptP)
	return eg
}
