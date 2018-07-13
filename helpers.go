package main

import (
	"github.com/go-macaron/session"
	"gopkg.in/macaron.v1"
)

func newPageData(ctx *macaron.Context, sess session.Store) {
	ctx.Data["URI"] = ctx.Req.RequestURI

	a := ctx.GetCookie("address")
	r := ctx.GetCookie("referral")

	if len(a) > 0 {
		u := &User{Address: a}
		db.First(u, u)
		if u.ID == 0 {
			if len(r) > 0 {
				u.Referral = r
			}
			db.Create(u)
		}
		ctx.Data["User"] = u
	}
}
