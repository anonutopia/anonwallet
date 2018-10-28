package main

import (
	"fmt"
	"log"

	"github.com/go-macaron/session"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/macaron.v1"
)

func homeView(ctx *macaron.Context) {
	ctx.HTML(200, "home")
}

func settingsView(ctx *macaron.Context) {
	user := ctx.Data["User"].(*User)
	var balance uint64

	abr, err := wnc.AssetsBalance(user.Address, "4zbprK67hsa732oSGLB6HzE8Yfdj3BcTcehCeTA1G5Lf")
	if err != nil {
		log.Printf("wnc.AssetsBalance error: %s", err)
		balance = 0
	} else {
		log.Println(abr.Balance)
		balance = uint64(abr.Balance)
	}

	if balance >= satInBtc {
		ctx.Data["UserHasAnotes"] = true
	} else {
		ctx.Data["UserHasAnotes"] = false
	}

	ctx.HTML(200, "settings")
}

func exchangeView(ctx *macaron.Context) {
	prices, err := pc.DoRequest()
	anote.loadState()
	if err != nil {
		ctx.Data["PriceWav"] = fmt.Sprintf("%.8f", float64(0))
		ctx.Data["PriceBtc"] = fmt.Sprintf("%.8f", float64(0))
		ctx.Data["PriceEth"] = fmt.Sprintf("%.8f", float64(0))
		ctx.Data["PriceEur"] = fmt.Sprintf("%.8f", float64(anote.Price)/float64(satInBtc))
	} else {
		ctx.Data["PriceWav"] = fmt.Sprintf("%.8f", prices.WAVES*(float64(anote.Price)/float64(satInBtc)))
		ctx.Data["PriceBtc"] = fmt.Sprintf("%.8f", prices.BTC*(float64(anote.Price)/float64(satInBtc)))
		ctx.Data["PriceEth"] = fmt.Sprintf("%.8f", prices.ETH*(float64(anote.Price)/float64(satInBtc)))
		ctx.Data["PriceEur"] = fmt.Sprintf("%.8f", float64(anote.Price)/float64(satInBtc))
	}
	ctx.HTML(200, "exchange")
}

func profitView(ctx *macaron.Context) {
	ctx.HTML(200, "profit")
}

func signInView(ctx *macaron.Context) {
	ctx.HTMLSet(200, "login", "signin")
}

func signInPostView(ctx *macaron.Context, siForm SignInForm, sess session.Store) {
	success := &JsonResponse{Success: false}

	user := ctx.Data["User"].(*User)

	if len(user.PasswordHash) == 0 {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(siForm.Password), 8)
		if err == nil {
			user.PasswordHash = string(hashedPassword)
			db.Save(user)
		}
	}

	err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(siForm.Password))
	if err == nil {
		success.Success = true
		sess.Set("userID", user.ID)
	}

	ctx.JSON(200, success)
}

func signOutView(ctx *macaron.Context, sess session.Store) {
	sess.Delete("userID")
	ctx.Redirect("/")
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

func signUpPostView(ctx *macaron.Context, siForm SignInForm, sess session.Store) {
	success := &JsonResponse{Success: true}

	user := ctx.Data["User"].(*User)
	log.Println(user)
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(siForm.Password), 8)
	if err == nil {
		user.PasswordHash = string(hashedPassword)
		db.Save(user)
	} else {
		success.Success = false
	}

	ctx.JSON(200, success)
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
	} else {
		err := sendWelcomeEmail(user)
		if err != nil {
			log.Printf("error sending email: %s", err)
		}
	}

	ctx.JSON(status, success)
}

func verifyView(ctx *macaron.Context, f *session.Flash, sess session.Store) {
	uid, err := decrypt([]byte(conf.DbPass[:16]), ctx.Params("uid"))
	if err != nil {
		return
	}

	log.Println(uid)
	u := &User{Address: uid}
	db.First(u, u)
	u.EmailVerified = true
	db.Save(u)

	sess.Set("userID", u.ID)
	ctx.SetCookie("address", u.Address, 1<<31-1)

	var balance uint64
	abr, err := wnc.AssetsBalance(u.Address, "4zbprK67hsa732oSGLB6HzE8Yfdj3BcTcehCeTA1G5Lf")
	if err != nil {
		balance = 0
	} else {
		balance = uint64(abr.Balance)
	}

	log.Printf("balance: %d", balance)

	f.Success("You have successfully verified your email address.")

	if balance >= (1 * satInBtc) {
		applicant := &Badge{Name: "applicant"}
		db.First(applicant, applicant)
		db.Model(u).Association("Badges").Append(applicant)

		f.Success("You have successfully applied for Anonutopia citizenship.")
	}

	if balance >= (1000 * satInBtc) {
		citizen := &Badge{Name: "citizen"}
		db.First(citizen, citizen)
		db.Model(u).Association("Badges").Append(citizen)
	}

	if balance >= (10000 * satInBtc) {
		founder := &Badge{Name: "funder"}
		db.First(founder, founder)
		db.Model(u).Association("Badges").Append(founder)
	}

	if balance >= (100000 * satInBtc) {
		pioneer := &Badge{Name: "pioneer"}
		db.First(pioneer, pioneer)
		db.Model(u).Association("Badges").Append(pioneer)
	}

	ctx.Redirect("/settings/")
}
