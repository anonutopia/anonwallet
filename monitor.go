package main

import (
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
		amountNew := balance - u.BitcoinBalanceProcessed
		if amountNew > 100000 {
			u.BitcoinBalanceNew = amountNew
			db.Save(u)
		}
	}
}

func (b *BitcoinAddressMonitor) checkAddressesRequest(address string) int {
	if len(address) == 0 {
		return 0
	}

	cl := http.Client{}

	url := fmt.Sprintf("https://blockchain.info/q/addressbalance/%s?confirmations=1", address)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Printf("[BitcoinAddressMonitor.checkAdressesRequest] request err %s", err)
		return 0
	}

	res, err := cl.Do(req)
	if err != nil {
		log.Printf("[BitcoinAddressMonitor.checkAdressesRequest] request do err %s", err)
		return 0
	}
	body, _ := ioutil.ReadAll(res.Body)

	balance, err := strconv.Atoi(string(body))
	if err == nil {
		return balance
	} else {
		log.Printf("[BitcoinAddressMonitor.checkAdressesRequest] strconv err: %s", err)
	}

	return 0
}

func initMonitor() *BitcoinAddressMonitor {
	bam := &BitcoinAddressMonitor{}
	bam.start()
	return bam
}
