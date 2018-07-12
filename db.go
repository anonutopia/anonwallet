package main

import (
	"log"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

func initDb() *gorm.DB {
	db, err := gorm.Open("sqlite3", DBNAME)

	if err != nil {
		log.Printf("[initDb] error: %s", err)
	}

	db.DB()
	db.DB().Ping()
	db.DB().SetMaxIdleConns(10)
	db.DB().SetMaxOpenConns(100)
	db.LogMode(conf.Debug)

	db.AutoMigrate(&User{})

	return db
}
