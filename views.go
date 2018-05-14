package main

import (
	"gopkg.in/macaron.v1"
)

func homeView(ctx *macaron.Context) {
	ctx.HTML(200, "home")
}
