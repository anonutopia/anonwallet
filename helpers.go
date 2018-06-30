package main

import (
	"github.com/go-macaron/session"
	"gopkg.in/macaron.v1"
)

func newPageData(ctx *macaron.Context, sess session.Store) {
	// ctx.Data["ProjectName"] = PROJECT_NAME
	ctx.Data["URI"] = ctx.Req.RequestURI
	// ctx.Data["ContractAddress"] = conf.ContractAddress
	ctx.Data["NetworkVersion"] = conf.EthNetwork

	// abi, err := ioutil.ReadFile("PonziCoin.abi")
	// if err == nil {
	// 	ctx.Data["ContractAbi"] = template.JS(string(abi))
	// }
}
