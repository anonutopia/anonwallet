package main

import (
	"github.com/jinzhu/gorm"
)

const (
	DBNAME = "wallet.db"
)

type User struct {
	gorm.Model
	Nickname string `sql:"size:255"`
	Email    string `sql:"size:255"`
	Address  string `sql:"size:255"`
	Seed     string `sql:"size:255"`
}
