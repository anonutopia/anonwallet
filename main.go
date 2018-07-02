package main

import (
	"fmt"
	"log"
	"net/http"

	"gopkg.in/macaron.v1"
)

var m *macaron.Macaron

var conf *config

func main() {
	m = initMacaron()

	conf = initConfig()

	m.Get("/", newPageData, homeView)
	m.Get("/profile/", newPageData, profileView)

	// m.Run()
	log.Println("Server is running...")
	http.ListenAndServe(fmt.Sprintf("0.0.0.0:%d", conf.Port), m)
}
