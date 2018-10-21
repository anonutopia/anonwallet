package main

import "log"

// const satInBtc = uint64(100000000)

const priceFactorLimit = uint64(0.0001 * float64(satInBtc))

type Anote struct {
	Price            uint64
	PriceFactor      uint64
	TierPrice        uint64
	TierPriceFactor  uint64
	BudgetWav        uint64
	BudgetBtc        uint64
	BudgetEth        uint64
	GatewayProfitBtc uint64
	GatewayProfitEth uint64
}

func (a *Anote) issueAmount(investment int, assetID string) int {
	p, err := pc.DoRequest()
	amount := int(0)
	if err == nil {
		var cryptoPrice float64

		if len(assetID) == 0 {
			cryptoPrice = p.WAVES
		} else if assetID == "7xHHNP8h6FrbP5jYZunYWgGn2KFSBiWcVaZWe644crjs" {
			cryptoPrice = p.BTC
		} else if assetID == "4fJ42MSLPXk9zwjfCdzXdUDAH8zQFCBdBz4sFSWZZY53" {
			cryptoPrice = p.ETH
		} else {
			return amount
		}

		for investment > 10 {
			log.Printf("anote: %d %d %d %d", a.Price, a.PriceFactor, a.TierPrice, a.TierPriceFactor)
			log.Printf("cryptoPrice: %f", cryptoPrice)
			log.Printf("investment: %d", investment)

			tierAmount := uint64(float64(investment) / cryptoPrice / float64(a.Price) * float64(satInBtc))

			log.Printf("tierAmount: %d", tierAmount)

			if tierAmount > a.TierPrice {
				tierAmount = a.TierPrice
			}

			log.Printf("tierAmount: %d", tierAmount)

			tierInvestment := int(float64(tierAmount) * float64(a.Price) * cryptoPrice / float64(satInBtc))

			log.Printf("tierInvestment: %d", tierInvestment)

			amount = amount + int(tierAmount)

			log.Printf("amount: %d", amount)

			investment = investment - tierInvestment

			log.Printf("investment: %d", investment)

			a.TierPrice = a.TierPrice - tierAmount
			a.TierPriceFactor = a.TierPriceFactor - tierAmount

			log.Printf("anote: %d %d %d %d", a.Price, a.PriceFactor, a.TierPrice, a.TierPriceFactor)

			if a.TierPrice == 0 {
				a.TierPrice = 1000 * satInBtc
				a.Price = a.Price + a.PriceFactor
			}

			if a.TierPriceFactor == 0 {
				a.TierPriceFactor = 1000000 * satInBtc
				if a.PriceFactor > priceFactorLimit {
					a.PriceFactor = a.PriceFactor - priceFactorLimit
				}
			}

			a.saveState()

			log.Printf("anote: %d %d %d %d", a.Price, a.PriceFactor, a.TierPrice, a.TierPriceFactor)
		}
	} else {
		log.Printf("[Anote.issueAmount] error pc.DoRequest: %s", err)
	}

	return amount
}

func (a *Anote) saveState() {
	ksip := &KeyValue{Key: "anotePrice"}
	db.FirstOrCreate(ksip, ksip)
	ksip.Value = a.Price
	db.Save(ksip)

	ksipf := &KeyValue{Key: "anotePriceFactor"}
	db.FirstOrCreate(ksipf, ksipf)
	ksipf.Value = a.PriceFactor
	db.Save(ksipf)

	ksitp := &KeyValue{Key: "anoteTierPrice"}
	db.FirstOrCreate(ksitp, ksitp)
	ksitp.Value = a.TierPrice
	db.Save(ksitp)

	ksitpf := &KeyValue{Key: "anoteTierPriceFactor"}
	db.FirstOrCreate(ksitpf, ksitpf)
	ksitpf.Value = a.TierPriceFactor
	db.Save(ksitpf)

	ksibw := &KeyValue{Key: "anoteBudgetWav"}
	db.FirstOrCreate(ksibw, ksibw)
	ksibw.Value = a.BudgetWav
	db.Save(ksibw)

	ksibb := &KeyValue{Key: "anoteBudgetBtc"}
	db.FirstOrCreate(ksibb, ksibb)
	ksibb.Value = a.BudgetBtc
	db.Save(ksibb)

	ksibe := &KeyValue{Key: "anoteBudgetEth"}
	db.FirstOrCreate(ksibe, ksibe)
	ksibe.Value = a.BudgetEth
	db.Save(ksibe)

	ksigpb := &KeyValue{Key: "anoteGatewayProfitBtc"}
	db.FirstOrCreate(ksigpb, ksigpb)
	ksigpb.Value = a.GatewayProfitBtc
	db.Save(ksigpb)

	ksigpd := &KeyValue{Key: "anoteGatewayProfitEth"}
	db.FirstOrCreate(ksigpd, ksigpd)
	ksigpd.Value = a.GatewayProfitEth
	db.Save(ksigpd)
}

func (a *Anote) loadState() {
	ksip := &KeyValue{Key: "anotePrice"}
	db.FirstOrCreate(ksip, ksip)

	if ksip.Value > 0 {
		a.Price = ksip.Value
	} else {
		ksip.Value = a.Price
		db.Save(ksip)
	}

	ksipf := &KeyValue{Key: "anotePriceFactor"}
	db.FirstOrCreate(ksipf, ksipf)

	if ksipf.Value > 0 {
		a.PriceFactor = ksipf.Value
	} else {
		ksipf.Value = a.PriceFactor
		db.Save(ksipf)
	}

	ksitp := &KeyValue{Key: "anoteTierPrice"}
	db.FirstOrCreate(ksitp, ksitp)

	if ksitp.Value > 0 {
		a.TierPrice = ksitp.Value
	} else {
		ksitp.Value = a.TierPrice
		db.Save(ksitp)
	}

	ksitpf := &KeyValue{Key: "anoteTierPriceFactor"}
	db.FirstOrCreate(ksitpf, ksitpf)

	if ksitpf.Value > 0 {
		a.TierPriceFactor = ksitpf.Value
	} else {
		ksitpf.Value = a.TierPriceFactor
		db.Save(ksitpf)
	}

	ksibw := &KeyValue{Key: "anoteBudgetWav"}
	db.FirstOrCreate(ksibw, ksibw)

	if ksibw.Value > 0 {
		a.BudgetWav = ksibw.Value
	} else {
		ksibw.Value = a.BudgetWav
		db.Save(ksibw)
	}

	ksibb := &KeyValue{Key: "anoteBudgetBtc"}
	db.FirstOrCreate(ksibb, ksibb)

	if ksibb.Value > 0 {
		a.BudgetBtc = ksibb.Value
	} else {
		ksibb.Value = a.BudgetBtc
		db.Save(ksibb)
	}

	ksibd := &KeyValue{Key: "anoteBudgetEth"}
	db.FirstOrCreate(ksibd, ksibd)

	if ksibd.Value > 0 {
		a.BudgetEth = ksibd.Value
	} else {
		ksibd.Value = a.BudgetEth
		db.Save(ksibd)
	}

	ksigpb := &KeyValue{Key: "anoteGatewayProfitBtc"}
	db.FirstOrCreate(ksigpb, ksigpb)

	if ksigpb.Value > 0 {
		a.GatewayProfitBtc = ksigpb.Value
	} else {
		ksigpb.Value = a.GatewayProfitBtc
		db.Save(ksigpb)
	}

	ksigpd := &KeyValue{Key: "anoteGatewayProfitEth"}
	db.FirstOrCreate(ksigpd, ksigpd)

	if ksigpd.Value > 0 {
		a.GatewayProfitEth = ksigpd.Value
	} else {
		ksigpd.Value = a.GatewayProfitEth
		db.Save(ksigpd)
	}
}

func initAnote() *Anote {
	anote := &Anote{
		Price:            uint64(0.01 * float64(satInBtc)),
		PriceFactor:      uint64(0.0021 * float64(satInBtc)),
		TierPrice:        1000 * satInBtc,
		TierPriceFactor:  1000000 * satInBtc,
		BudgetWav:        0,
		BudgetBtc:        0,
		BudgetEth:        0,
		GatewayProfitBtc: 0,
		GatewayProfitEth: 0}

	anote.loadState()

	return anote
}
