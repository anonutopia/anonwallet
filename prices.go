package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

const url = "https://min-api.cryptocompare.com/data/price?fsym=EUR&tsyms=WAVES,BTC,ETH"

type Prices struct {
	WAVES float64 `json:"WAVES"`
	BTC   float64 `json:"BTC"`
	ETH   float64 `json:"ETH"`
}

type PriceClient struct {
	Url string
}

func (w *PriceClient) DoRequest() (*Prices, error) {
	p := &Prices{}
	cl := http.Client{}

	var req *http.Request
	var err error

	req, err = http.NewRequest(http.MethodGet, url, nil)

	req.Header.Set("Content-Type", "application/json")

	if err != nil {
		return nil, err
	}

	res, err := cl.Do(req)

	if err == nil {
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return nil, err
		}
		if res.StatusCode != 200 {
			log.Printf("[PriceClient.DoRequest] Error, body: %s", string(body))
		}
		json.Unmarshal(body, p)
	} else {
		return nil, err
	}

	return p, nil
}

func initPriceClient() *PriceClient {
	pc := &PriceClient{}
	return pc
}
