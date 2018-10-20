package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"
)

type BitcoinAddressMonitor struct {
}

func (b *BitcoinAddressMonitor) start() {
	go func() {
		for {
			time.Sleep(time.Second * 5)

			b.checkAddresses()
		}
	}()
}

func (b *BitcoinAddressMonitor) checkAddresses() {
	var users []*User
	db.Find(&users)

	for _, u := range users {
		balance := b.checkAddressesRequest(u.BitcoinAddr)
		if balance > 0 {
			u.BitcoinBalanceNew = balance

			ua := &UsedAddress{Address: u.BitcoinAddr, Type: 1, UserID: int(u.ID), Balance: uint64(balance)}
			db.Create(ua)

			var err error
			u.BitcoinAddr, err = bg.getAddress()
			if err != nil {
				log.Printf("Error in bg.getAddress: %s", err)
			}

			db.Save(u)
		}
	}
}

func (b *BitcoinAddressMonitor) checkAddressesRequest(address string) int {
	if len(address) == 0 {
		return 0
	}

	balance, err := bg.getBalance(address)
	if err != nil {
		log.Printf("Error in BitcoinAddressMonitor.checkAddressRequest: %s", err)
		return 0
	}

	return int(balance * 100000000)
}

func initBaMonitor() *BitcoinAddressMonitor {
	bam := &BitcoinAddressMonitor{}
	bam.start()
	return bam
}

type EthereumAddressMonitor struct {
}

func (e *EthereumAddressMonitor) start() {
	go func() {
		for {
			time.Sleep(time.Second * 5)

			e.checkAddresses()
		}
	}()
}

func (e *EthereumAddressMonitor) checkAddresses() {
	var users []*User
	db.Find(&users)

	for _, u := range users {
		balance := e.checkAddressesRequest(u.EtherAddr)
		if balance > 0 {
			u.EtherBalanceNew = balance

			ua := &UsedAddress{Address: u.EtherAddr, Type: 2, UserID: int(u.ID), Balance: uint64(balance)}
			db.Create(ua)

			var err error
			u.EtherAddr, err = eg.getAddress(uint32(u.ID))
			if err != nil {
				log.Printf("Error in eg.getAddress: %s", err)
			}

			db.Save(u)
		}
	}
}

func (e *EthereumAddressMonitor) checkAddressesRequest(address string) int {
	if len(address) == 0 {
		return 0
	}

	cl := http.Client{}

	url := fmt.Sprintf("https://api.etherscan.io/api?module=account&action=balance&tag=latest&address=%s", address)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Printf("[EthereumAddressMonitor.checkAdressesRequest] request err %s", err)
		return 0
	}

	res, err := cl.Do(req)
	if err != nil {
		log.Printf("[EthereumAddressMonitor.checkAdressesRequest] request do err %s", err)
		return 0
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Printf("[EthereumAddressMonitor.checkAdressesRequest] ioutil.ReadAll err %s", err)
		return 0
	}

	b := &EthBalance{}

	json.Unmarshal(body, b)

	balance, err := strconv.Atoi(b.Result)
	if err != nil {
		return 0
	}

	return balance / 10000000000
}

func initEaMonitor() *EthereumAddressMonitor {
	eam := &EthereumAddressMonitor{}
	eam.start()
	return eam
}

type EthBalance struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Result  string `json:"result"`
}
