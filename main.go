package main

import (
	"log"

	"gopkg.in/macaron.v1"
)

var m *macaron.Macaron

var conf *config

func main() {
	m = initMacaron()

	conf = initConfig()

	m.Get("/", newPageData, homeView)

	// m.Run()
	log.Println("Server is running...")
	// http.ListenAndServe("0.0.0.0:4001", m)
	m.Run()
}
