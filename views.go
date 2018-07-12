package main

import (
	"gopkg.in/macaron.v1"
)

func homeView(ctx *macaron.Context) {
	ctx.HTML(200, "home")
}

func profileView(ctx *macaron.Context) {
	ctx.HTML(200, "profile")
}

func exchangeView(ctx *macaron.Context) {
	ctx.HTML(200, "exchange")
}

func profitView(ctx *macaron.Context) {
	ctx.HTML(200, "profit")
}

func signInView(ctx *macaron.Context) {
	ctx.HTMLSet(200, "login", "signin")
}

func signUpView(ctx *macaron.Context) {
	ctx.HTMLSet(200, "login", "signup")
}
