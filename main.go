package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/anonutopia/gowaves"
	"github.com/go-macaron/binding"
	"github.com/jinzhu/gorm"
	"gopkg.in/macaron.v1"
)

var m *macaron.Macaron

var conf *config

var db *gorm.DB

var bg *BitcoinGenerator

var eg *EthereumGenerator

var bam *BitcoinAddressMonitor

var eam *EthereumAddressMonitor

var loc map[string]string

var pc *PriceClient

var wnc *gowaves.WavesNodeClient

var anote *Anote

func main() {
	m = initMacaron()

	conf = initConfig()

	db = initDb()

	bg = initBtcGen()

	eg = initEthGen()

	bam = initBaMonitor()

	eam = initEaMonitor()

	pc = initPriceClient()

	wnc = initWaves()

	anote = initAnote()

	m.Get("/", newPageData, loginRequired, homeView)
	m.Get("/settings/", newPageData, loginRequired, settingsView)
	m.Get("/exchange/", newPageData, loginRequired, exchangeView)
	m.Get("/profit/", newPageData, loginRequired, profitView)
	m.Get("/sign-in/", newPageData, signInView)
	m.Get("/sign-out/", newPageData, loginRequired, signOutView)
	m.Get("/sign-up/", newPageData, signUpView)
	m.Get("/sign-up-new/", newPageData, signUpNewView)
	m.Get("/sign-up-import/", newPageData, signUpImportView)
	m.Get("/locales.json", newPageData, localesjsView)
	m.Get("/verify/:uid", newPageData, verifyView)

	m.Post("/apply/", binding.Bind(ApplyForm{}), newPageData, loginRequired, applyView)
	m.Post("/sign-in/", binding.Bind(SignInForm{}), newPageData, signInPostView)
	m.Post("/sign-up/", binding.Bind(SignInForm{}), newPageData, signUpPostView)

	// m.Run()
	addr := fmt.Sprintf("0.0.0.0:%d", conf.Port)
	log.Printf("Server is running on: %s", addr)
	http.ListenAndServe(addr, m)
}
