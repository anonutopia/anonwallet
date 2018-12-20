package main

import (
	"github.com/go-macaron/binding"
	macaron "gopkg.in/macaron.v1"
)

type ProfileForm struct {
	Nickname string `binding:"Required"`
	Email    string `binding:"Email"`
	Country  string `binding:"Required"`
	City     string `binding:"Required"`
}

type SignUpForm struct {
	Password string `binding:"Required"`
	Email    string `binding:"Required;Email"`
	Seed     string `binding:"Required"`
}

func (lf SignUpForm) Error(ctx *macaron.Context, errs binding.Errors) {
	ctx.Data["Errors"] = errs
}

type FacebookAwardForm struct {
	FbLink string `binding:"Required"`
}

func (faf FacebookAwardForm) Error(ctx *macaron.Context, errs binding.Errors) {
	ctx.Data["Errors"] = errs
}
