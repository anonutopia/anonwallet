package main

import (
	"fmt"

	"gopkg.in/macaron.v1"
)

func homeView(ctx *macaron.Context) {
	ctx.HTML(200, "home")
}

func settingsView(ctx *macaron.Context) {
	ctx.HTML(200, "settings")
}

func exchangeView(ctx *macaron.Context) {
	prices, err := pc.DoRequest()
	if err != nil {
		ctx.Data["PriceWav"] = fmt.Sprintf("%.8f", float64(0))
		ctx.Data["PriceBtc"] = fmt.Sprintf("%.8f", float64(0))
		ctx.Data["PriceEth"] = fmt.Sprintf("%.8f", float64(0))
	} else {
		ctx.Data["PriceWav"] = fmt.Sprintf("%.8f", prices.WAVES/100)
		ctx.Data["PriceBtc"] = fmt.Sprintf("%.8f", prices.BTC/100)
		ctx.Data["PriceEth"] = fmt.Sprintf("%.8f", prices.ETH/100)
	}
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

func signUpNewView(ctx *macaron.Context) {
	ctx.HTMLSet(200, "login", "signupnew")
}

func signUpImportView(ctx *macaron.Context) {
	ctx.HTMLSet(200, "login", "signupimport")
}

func localesjsView(ctx *macaron.Context) {
	loc = initLocale(ctx)
	ctx.JSON(200, &loc)
}

func applyView(ctx *macaron.Context, af ApplyForm) {
	success := &JsonResponse{Success: true}
	status := 200
	user := ctx.Data["User"].(*User)

	user.Nickname = af.Nickname
	user.Email = af.Email
	user.Country = af.Country
	user.City = af.City

	err := db.Save(user)

	if err.RowsAffected == 0 {
		success.Success = false
		success.Message = err.Error.Error()
		status = 400
	}

	ctx.JSON(status, success)
}
