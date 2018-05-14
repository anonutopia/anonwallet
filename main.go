package main

import (
	"log"

	"gopkg.in/macaron.v1"
)

var m *macaron.Macaron

func main() {
	m = initMacaron()

	m.Get("/", homeView)

	// m.Run()
	log.Println("Server is running...")
	// http.ListenAndServe("0.0.0.0:4001", m)
	m.Run()
}
