package main

import (
	"log"

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
		if u.ID != 0 && len(u.BitcoinAddr) == 0 {
			var err error
			u.BitcoinAddr, err = bg.getAddress()
			if err != nil {
				log.Printf("Error in bg.getAddress: %s", err)
			}
			db.Save(u)
		}
		if u.ID != 0 && len(u.EtherAddr) == 0 {
			var err error
			u.EtherAddr, err = eg.getAddress(uint32(u.ID))
			if err != nil {
				log.Printf("Error in eg.getAddress: %s", err)
			}
			db.Save(u)
		}
		ctx.Data["User"] = u
		sess.Set("userID", u.ID)
	}
}

type JsonResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

func loginRequired(ctx *macaron.Context, sess session.Store) {
	userID := sess.Get("userID")
	if userID == nil {
		ctx.Redirect("/sign-in/")
	}
}
