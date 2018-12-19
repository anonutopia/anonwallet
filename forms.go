package main

import (
	"github.com/go-macaron/binding"
	macaron "gopkg.in/macaron.v1"
)

type ApplyForm struct {
	Nickname string `binding:"Required"`
	Email    string `binding:"Email"`
	Country  string `binding:"Required"`
	City     string `binding:"Required"`
}

type SignInForm struct {
	Password string `binding:"Required"`
	Seed     string
}

type FacebookAwardForm struct {
	FbLink string `binding:"Required"`
}

func (faf FacebookAwardForm) Error(ctx *macaron.Context, errs binding.Errors) {
	ctx.Data["Errors"] = errs
}
