package main

import (
	"github.com/go-macaron/session"
	"gopkg.in/macaron.v1"
)

func newPageData(ctx *macaron.Context, sess session.Store) {
	ctx.Data["URI"] = ctx.Req.RequestURI

	a := ctx.GetCookie("address")

	if len(a) > 0 {
		u := &User{Address: a}
		db.FirstOrCreate(u)
		ctx.Data["User"] = u
	}
}
