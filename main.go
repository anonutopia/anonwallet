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

func main() {
	m = initMacaron()

	conf = initConfig()

	db = initDb()

	m.Get("/", newPageData, homeView)
	m.Get("/profile/", newPageData, profileView)
	m.Get("/exchange/", newPageData, exchangeView)
	m.Get("/profit/", newPageData, profitView)
	m.Get("/sign-in/", newPageData, signInView)

	// m.Run()
	log.Println("Server is running...")
	http.ListenAndServe(fmt.Sprintf("0.0.0.0:%d", conf.Port), m)
}
