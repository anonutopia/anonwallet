package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/jinzhu/gorm"
	"gopkg.in/macaron.v1"
)

var m *macaron.Macaron

var conf *config

var db *gorm.DB

var bg *BitcoinGenerator

var bam *BitcoinAddressMonitor

func main() {
	m = initMacaron()

	conf = initConfig()

	db = initDb()

	bg = initBtcGen()

	bam = initMonitor()

	m.Get("/", newPageData, homeView)
	m.Get("/profile/", newPageData, profileView)
	m.Get("/exchange/", newPageData, exchangeView)
	m.Get("/profit/", newPageData, profitView)
	m.Get("/sign-in/", newPageData, signInView)
	m.Get("/sign-up/", newPageData, signUpView)
	m.Get("/sign-up-new/", newPageData, signUpNewView)
	m.Get("/sign-up-import/", newPageData, signUpImportView)

	// m.Run()
	addr := fmt.Sprintf("0.0.0.0:%d", conf.Port)
	log.Printf("Server is running on: %s", addr)
	http.ListenAndServe(addr, m)
}
