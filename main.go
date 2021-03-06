package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/anonutopia/gowaves"
	"github.com/go-macaron/binding"
	"github.com/jinzhu/gorm"
	macaron "gopkg.in/macaron.v1"
	tgbotapi "gopkg.in/telegram-bot-api.v4"
)

var m *macaron.Macaron

var conf *config

var bot *tgbotapi.BotAPI

var db *gorm.DB

var bg *BitcoinGenerator

var eg *EthereumGenerator

var bam *BitcoinAddressMonitor

var eam *EthereumAddressMonitor

var ubm *UserBalanceMonitor

var loc map[string]string

var pc *PriceClient

var wnc *gowaves.WavesNodeClient

var anote *Anote

func main() {
	m = initMacaron()
	conf = initConfig()
	bot = initBot()
	db = initDb()
	bg = initBtcGen()
	eg = initEthGen()
	bam = initBaMonitor()
	eam = initEaMonitor()
	pc = initPriceClient()
	wnc = initWaves()
	anote = initAnote()
	ubm = initUserBalanceMonitor()

	logTelegram("Successfully started - anonwallet.")

	m.Get("/", newPageData, loginRequired, homeView)
	m.Get("/settings/", newPageData, loginRequired, settingsView)
	m.Get("/exchange/", newPageData, loginRequired, exchangeView)
	m.Get("/profit/", newPageData, loginRequired, profitView)
	m.Get("/sign-out/", newPageData, loginRequired, signOutView)
	m.Get("/sign-up/", newPageData, signUpView)
	m.Get("/sign-in/", newPageData, signInView)
	m.Get("/sign-in-clean/", newPageData, signInOldCleanView)
	m.Get("/locales.json", newPageData, localesjsView)
	m.Get("/verify/:uid", newPageData, verifyView)
	m.Get("/password-reset/", newPageData, passwordResetView)
	m.Get("/password-reset/:uid", newPageData, passwordResetSignInView)
	m.Get("/password-reset-finish/", newPageData, loginRequired, passwordResetFinishView)

	m.Post("/apply/", binding.Bind(ProfileForm{}), newPageData, loginRequired, applyView)
	m.Post("/settings/", binding.Bind(FacebookAwardForm{}), newPageData, initFbPostView)
	m.Post("/sign-up/", binding.Bind(SignUpForm{}), newPageData, signUpPostView)
	m.Post("/sign-in/", binding.Bind(SignInForm{}), newPageData, signInPostView)
	m.Post("/sign-in-clean/", binding.Bind(SignUpForm{}), newPageData, signInOldCleanPostView)
	m.Post("/password-reset/", binding.Bind(SignInOldForm{}), newPageData, passwordResetPostView)
	m.Post("/password-reset-finish/", binding.Bind(SignInForm{}), newPageData, loginRequired, passwordResetFinishPostView)

	addr := fmt.Sprintf("0.0.0.0:%d", conf.Port)
	log.Printf("Server is running on: %s", addr)
	http.ListenAndServe(addr, m)
}
