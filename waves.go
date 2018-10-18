package main

import (
	"github.com/anonutopia/gowaves"
)

func initWaves() *gowaves.WavesNodeClient {
	wnc := &gowaves.WavesNodeClient{
		Host:   "localhost",
		Port:   6869,
		ApiKey: conf.WavesNodeApiKey,
	}

	return wnc
}
