package main

import (
	"fmt"
	"log"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

func initDb() *gorm.DB {
	db, err := gorm.Open("postgres", fmt.Sprintf("host=localhost port=5432 user=%s dbname=%s password=%s sslmode=disable", conf.DbUser, conf.DbName, conf.DbPass))

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
