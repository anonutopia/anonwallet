package main

import (
	"fmt"

	"github.com/jinzhu/gorm"
)

const (
	DBNAME = "wallet.db"
)

type User struct {
	gorm.Model
	Nickname  string `sql:"size:255"`
	Email     string `sql:"size:255"`
	Address   string `sql:"size:255"`
	Seed      string `sql:"size:255"`
	ProfitEth uint64
	ProfitWav uint64
	ProfitBtc uint64
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
