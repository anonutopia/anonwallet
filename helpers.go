package main

import (
	"html/template"
	"io/ioutil"

	"github.com/go-macaron/session"
	"gopkg.in/macaron.v1"
)

func newPageData(ctx *macaron.Context, sess session.Store) {
	ctx.Data["URI"] = ctx.Req.RequestURI
	ctx.Data["AnoteAddress"] = conf.AnoteAddress
	ctx.Data["NetworkVersion"] = conf.EthNetwork

	abi, err := ioutil.ReadFile("abi/ANT.abi")
	if err == nil {
		ctx.Data["AnoteAbi"] = template.JS(string(abi))
	}
}
