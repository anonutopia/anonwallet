package main

import (
	"fmt"

	"github.com/jinzhu/gorm"
)

type User struct {
	gorm.Model
	Nickname               string `sql:"size:255"`
	Email                  string `sql:"size:255"`
	Address                string `sql:"size:255"`
	Referral               string `sql:"size:255"`
	ProfitEth              uint64
	ProfitWav              uint64
	ProfitBtc              uint64
	ProfitEthTotal         uint64
	ProfitWavTotal         uint64
	ProfitBtcTotal         uint64
	ReferralProfitEth      uint64
	ReferralProfitWav      uint64
	ReferralProfitBtc      uint64
	ReferralProfitEthTotal uint64
	ReferralProfitWavTotal uint64
	ReferralProfitBtcTotal uint64
}

func (u *User) ProfitWavString() string {
	return fmt.Sprintf("%.8f", float64(u.ProfitWav)/float64(100000000))
}

func (u *User) ProfitBtcString() string {
	return fmt.Sprintf("%.8f", float64(u.ProfitBtc)/float64(100000000))
}

func (u *User) ProfitEthString() string {
	return fmt.Sprintf("%.8f", float64(u.ProfitEth)/float64(100000000))
}

func (u *User) ReferralProfitWavString() string {
	return fmt.Sprintf("%.8f", float64(u.ReferralProfitWav)/float64(100000000))
}

func (u *User) ReferralProfitBtcString() string {
	return fmt.Sprintf("%.8f", float64(u.ReferralProfitBtc)/float64(100000000))
}

func (u *User) ReferralProfitEthString() string {
	return fmt.Sprintf("%.8f", float64(u.ReferralProfitEth)/float64(100000000))
}
