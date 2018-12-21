package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"reflect"
	"strings"
	"time"

	"github.com/anonutopia/gowaves"
	"github.com/btcsuite/golangcrypto/bcrypt"
	"github.com/go-macaron/session"
	macaron "gopkg.in/macaron.v1"
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
		logTelegram(fmt.Sprintf("wnc.AssetsBalance error: %s", err))
		balance = 0
	} else {
		balance = uint64(abr.Balance)
	}

	if balance >= satInBtc {
		ctx.Data["UserHasAnotes"] = true
	} else {
		ctx.Data["UserHasAnotes"] = false
	}

	ctx.Data["Referred"] = user.ReferredUsersVerifiedCount()

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
	user := ctx.Data["User"].(*User)
	prices, _ := pc.DoRequest()

	tp := float64(0)

	tp += float64(user.ProfitWav) / (prices.WAVES * float64(satInBtc))
	tp += float64(user.ProfitBtc) / (prices.BTC * float64(satInBtc))
	tp += float64(user.ProfitEth) / (prices.ETH * float64(satInBtc))
	tp += float64(user.ReferralProfitWav) / (prices.WAVES * float64(satInBtc))
	tp += float64(user.ReferralProfitBtc) / (prices.BTC * float64(satInBtc))
	tp += float64(user.ReferralProfitEth) / (prices.ETH * float64(satInBtc))

	ctx.Data["TotalProfit"] = fmt.Sprintf("%.2f", tp)

	ctx.HTML(200, "profit")
}

func signOutView(ctx *macaron.Context, sess session.Store) {
	sess.Delete("userID")
	ctx.Redirect("/")
}

func signUpView(ctx *macaron.Context) {
	ctx.Data["Title"] = "Apply for Citizenship | "

	ctx.HTMLSet(200, "login", "signup")
}

func signUpPostView(ctx *macaron.Context, suf SignUpForm, sess session.Store) {
	ctx.Data["Title"] = "Apply for Citizenship | "
	ctx.Data["SignUpForm"] = suf

	s := reflect.ValueOf(ctx.Data["Errors"])

	if s.Len() == 0 {
		u := &User{Email: suf.Email}
		db.First(u, u)
		if u.ID != 0 {
			ctx.Data["Errors"] = true
			ctx.Data["ErrorMsg"] = "This user already exists."

			ctx.Data["OldUserWarning"] = true
		} else if !validateEmailDomain(suf.Email) {
			ctx.Data["Errors"] = true
			ctx.Data["ErrorMsg"] = "Please use one of known email providers like Gmail."
		} else {
			r := ctx.GetCookie("referral")
			u.Address = suf.Address
			u.Nickname = u.Email
			hashedPassword, err := bcrypt.GenerateFromPassword([]byte(suf.Password), 8)
			for err != nil {
				time.Sleep(200 * time.Microsecond)
				hashedPassword, err = bcrypt.GenerateFromPassword([]byte(suf.Password), 8)
			}
			u.PasswordHash = string(hashedPassword)
			if len(r) > 0 {
				u.Referral = r
			}
			db.Create(u)

			err = sendInitWelcomeEmail(u, suf.Password, suf.Seed, "en-US")
			if err != nil {
				log.Printf("Error in send welcome email: %s", err)
				logTelegram(fmt.Sprintf("Error in send welcome email: %s", err))
			}

			sess.Set("userID", u.ID)

			ctx.Data["Finished"] = true
		}
	} else {
		ctx.Data["ErrorMsg"] = "Email address is required."
	}

	ctx.HTMLSet(200, "login", "signup")
}

func signInView(ctx *macaron.Context) {
	ctx.Data["Title"] = "Wallet Sign In | "

	ctx.HTMLSet(200, "login", "signin")
}

func signInPostView(ctx *macaron.Context, sif SignInForm, sess session.Store) {
	ctx.Data["Title"] = "Wallet Sign In | "
	ctx.Data["SignInForm"] = sif

	s := reflect.ValueOf(ctx.Data["Errors"])

	if s.Len() == 0 {
		user := &User{Address: sif.Address}
		db.First(user, user)

		err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(sif.Password))
		if err == nil {
			sess.Set("userID", user.ID)

			ctx.Data["Finished"] = true
		} else {
			ctx.Data["Errors"] = true
			ctx.Data["ErrorMsg"] = "Wrong password, please try again."
		}
	} else {
		ctx.Data["ErrorMsg"] = "Password is required."
	}

	ctx.HTMLSet(200, "login", "signin")
}

func signInOldView(ctx *macaron.Context) {
	ctx.Data["Title"] = "Old Wallet Sign In | "

	ctx.HTMLSet(200, "login", "signinold")
}

func signInOldPostView(ctx *macaron.Context, siof SignInOldForm, f *session.Flash) {
	ctx.Data["Title"] = "Old Wallet Sign In | "
	ctx.Data["SignInOldForm"] = siof

	s := reflect.ValueOf(ctx.Data["Errors"])

	if s.Len() == 0 {
		// user := ctx.Data["User"].(*User)
		if validateEmailDomain(siof.Email) {
			user := &User{Email: siof.Email}
			db.First(user, user)
			if user.ID == 0 {
				user := ctx.Data["User"].(*User)
				user.Email = siof.Email
				db.Save(user)
				f.Success("You have successfully signed in to your Anonutopia wallet.")
				ctx.Redirect("/")
			} else {
				ctx.Data["Errors"] = true
				ctx.Data["ErrorMsg"] = "This email address is already used."
			}
		} else {
			ctx.Data["Errors"] = true
			ctx.Data["ErrorMsg"] = "Please use one of known email providers like Gmail."
		}
		// user := &User{Address: sif.Address}
		// db.First(user, user)

		// err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(sif.Password))
		// if err == nil {
		// 	sess.Set("userID", user.ID)

		// 	ctx.Data["Finished"] = true
		// } else {
		// 	ctx.Data["Errors"] = true
		// 	ctx.Data["ErrorMsg"] = "Wrong password, please try again."
		// }
	} else {
		ctx.Data["ErrorMsg"] = "Email is required."
	}

	ctx.HTMLSet(200, "login", "signinold")
}

func signInOldCleanView(ctx *macaron.Context) {
	ctx.Data["Title"] = "Old Wallet Clean Sign In | "

	ctx.HTMLSet(200, "login", "signinoldclean")
}

func signInOldCleanPostView(ctx *macaron.Context, suf SignUpForm, sess session.Store) {
	ctx.Data["Title"] = "Old Wallet Clean Sign In | "
	ctx.Data["SignUpForm"] = suf

	s := reflect.ValueOf(ctx.Data["Errors"])

	if s.Len() == 0 {
		u := &User{Email: suf.Email}
		db.First(u, u)
		if u.ID != 0 {
			if u.Address == suf.Address {
				err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(suf.Password))
				log.Println(err)
				if err == nil {
					sess.Set("userID", u.ID)
					ctx.Data["Finished"] = true
				} else {
					ctx.Data["Errors"] = true
					ctx.Data["ErrorMsg"] = "Wrong password, please try again."
				}
			} else {
				ctx.Data["Errors"] = true
				ctx.Data["ErrorMsg"] = "Something is wrong with your seed (extra space or new line)."
			}
		} else {
			u := &User{Address: suf.Address}
			db.First(u, u)

			if u.ID != 0 {
				// sess.Set("userID", u.ID)
				// ctx.Data["Finished"] = true
			} else {
				ctx.Data["Errors"] = true
				ctx.Data["ErrorMsg"] = "The user with this email or address doesn't exist."
			}
		}
	} else {
		ctx.Data["ErrorMsg"] = "All fields are required."
	}

	ctx.HTMLSet(200, "login", "signinoldclean")
}

func localesjsView(ctx *macaron.Context) {
	loc = initLocale(ctx)
	ctx.JSON(200, &loc)
}

func applyView(ctx *macaron.Context, af ProfileForm) {
	success := &JsonResponse{Success: true}
	status := 200

	if validateEmailDomain(af.Email) {
		user := ctx.Data["User"].(*User)
		first := false

		user.Nickname = af.Nickname
		if len(user.Email) == 0 {
			user.Email = af.Email
			first = true
		}
		user.Country = af.Country
		user.City = af.City

		err := db.Save(user)

		if err.RowsAffected == 0 {
			success.Success = false
			success.Message = err.Error.Error()
			status = 400
		} else if first {
			err := sendWelcomeEmail(user, ctx.GetCookie("lang"))
			if err != nil {
				log.Printf("error sending email: %s", err)
				logTelegram(fmt.Sprintf("error sending email: %s", err))
			} else {
				applicant := &Badge{Name: "applicant"}
				db.First(applicant, applicant)
				db.Model(user).Association("Badges").Append(applicant)
			}
		}
	}

	ctx.JSON(status, success)
}

func verifyView(ctx *macaron.Context, f *session.Flash, sess session.Store) {
	uid, err := decrypt([]byte(conf.DbPass[:16]), ctx.Params("uid"))
	if err != nil {
		return
	}

	u := &User{Email: uid}
	db.First(u, u)

	if !u.EmailVerified {
		atr := &gowaves.AssetsTransferRequest{
			Amount:    100000000,
			AssetID:   "4zbprK67hsa732oSGLB6HzE8Yfdj3BcTcehCeTA1G5Lf",
			Fee:       100000,
			Recipient: uid,
			Sender:    conf.NodeAddress,
		}

		_, err := wnc.AssetsTransfer(atr)
		if err != nil {
			log.Printf("Error gowaves.AssetsTransfer %s", err)
			logTelegram(fmt.Sprintf("Error gowaves.AssetsTransfer %s", err))
		}

		u.EmailVerified = true
		db.Save(u)

		if len(u.Referral) > 0 {
			atr := &gowaves.AssetsTransferRequest{
				Amount:    20000000,
				AssetID:   "4zbprK67hsa732oSGLB6HzE8Yfdj3BcTcehCeTA1G5Lf",
				Fee:       100000,
				Recipient: u.Referral,
				Sender:    conf.NodeAddress,
			}

			wnc.AssetsTransfer(atr)

			r := &User{Address: u.Referral}
			db.First(r, r)
			count := r.ReferredUsersVerifiedCount()
			if count >= 10 {
				citizen := &Badge{Name: "citizen"}
				db.First(citizen, citizen)
				db.Model(r).Association("Badges").Append(citizen)
			}
		}
	}

	sess.Set("userID", u.ID)
	ctx.SetCookie("address", u.Address, 1<<31-1)

	var balance uint64
	abr, err := wnc.AssetsBalance(u.Address, "4zbprK67hsa732oSGLB6HzE8Yfdj3BcTcehCeTA1G5Lf")
	if err != nil {
		balance = 0
	} else {
		balance = uint64(abr.Balance)
	}

	f.Success("You have successfully verified your email address. We have sent you your 1 free anote.")

	citizenLimit := (10 * satInBtc / anote.Price) * satInBtc

	if balance >= citizenLimit {
		citizen := &Badge{Name: "citizen"}
		db.First(citizen, citizen)
		db.Model(u).Association("Badges").Append(citizen)
	}

	if balance >= (10000 * satInBtc) {
		founder := &Badge{Name: "founder"}
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

func initFbPostView(ctx *macaron.Context, faf FacebookAwardForm, f *session.Flash) {
	ctx.Data["ErrorMessage"] = "This field is required."
	ctx.Data["Form"] = faf
	user := ctx.Data["User"].(*User)

	if strings.Contains(faf.FbLink, "facebook.com/") {
		fbLinkSplit := strings.Split(faf.FbLink, "facebook.com/")
		sp := &SharePost{URI: fbLinkSplit[1]}
		db.First(sp, sp)

		if sp.ID == 0 {
			response, err := http.Get(faf.FbLink)
			if err != nil {
				ctx.Data["ErrorMessage"] = fmt.Sprintf("%e", err)
			} else {
				defer response.Body.Close()
				contents, err := ioutil.ReadAll(response.Body)
				if err != nil {
					ctx.Data["ErrorMessage"] = fmt.Sprintf("%e", err)
				} else {
					containsLink := strings.Contains(string(contents), "https://www.anonutopia.com") || strings.Contains(string(contents), "https%3A%2F%2Fwww.anonutopia.com")
					containsAddress := strings.Contains(string(contents), user.Address)
					timeLimitAllowed := false

					if user.LastFacebookAwardTime == nil {
						timeLimitAllowed = true
					} else if user.LastFacebookAwardTime.Add(24 * time.Hour).Before(time.Now()) {
						timeLimitAllowed = true
					}

					if containsLink && containsAddress && timeLimitAllowed && user.NextFacebookAward > 0 {
						atr := &gowaves.AssetsTransferRequest{
							Amount:    user.NextFacebookAward,
							AssetID:   "4zbprK67hsa732oSGLB6HzE8Yfdj3BcTcehCeTA1G5Lf",
							Fee:       100000,
							Recipient: user.Address,
							Sender:    conf.NodeAddress,
						}

						_, err := wnc.AssetsTransfer(atr)
						if err != nil {
							log.Printf("Error gowaves.AssetsTransfer %s", err)
							logTelegram(fmt.Sprintf("Error gowaves.AssetsTransfer %s", err))
						} else {
							f.Success("You have successfully received your Facebook share award.")

							now := time.Now()
							user.NextFacebookAward -= 20000000
							user.LastFacebookAwardTime = &now
							db.Save(user)

							db.Create(sp)

							ctx.Redirect("/settings/")
							return
						}
					} else if !containsLink || !containsAddress {
						ctx.Data["Errors"] = true
						ctx.Data["ErrorMessage"] = "Pasted URL doesn't containt your referral link."
					} else if !timeLimitAllowed {
						ctx.Data["Errors"] = true
						ctx.Data["ErrorMessage"] = "You have to wait at least 24 hours to do this again."
					} else if user.NextFacebookAward == 0 {
						ctx.Data["Errors"] = true
						ctx.Data["ErrorMessage"] = "You have received all Facebook share awards."
					}
				}
			}
		} else {
			ctx.Data["Errors"] = true
			ctx.Data["ErrorMessage"] = "This Facebook post has already been used."
		}
	} else {
		ctx.Data["Errors"] = true
		ctx.Data["ErrorMessage"] = "The link you provided is not a Facebook link."
	}

	var balance uint64

	abr, err := wnc.AssetsBalance(user.Address, "4zbprK67hsa732oSGLB6HzE8Yfdj3BcTcehCeTA1G5Lf")
	if err != nil {
		log.Printf("wnc.AssetsBalance error: %s", err)
		logTelegram(fmt.Sprintf("wnc.AssetsBalance error: %s", err))
		balance = 0
	} else {
		balance = uint64(abr.Balance)
	}

	if balance >= satInBtc {
		ctx.Data["UserHasAnotes"] = true
	} else {
		ctx.Data["UserHasAnotes"] = false
	}

	ctx.Data["Referred"] = user.ReferredUsersVerifiedCount()

	ctx.HTML(200, "settings")
}

func passwordResetView(ctx *macaron.Context) {
	ctx.Data["Title"] = "Password Reset | "

	ctx.HTMLSet(200, "login", "passwordreset")
}

func passwordResetPostView(ctx *macaron.Context, form SignInOldForm) {
	ctx.Data["Title"] = "Password Reset | "
	ctx.Data["Form"] = form

	uid, err := encrypt([]byte(conf.DbPass[:16]), form.Email)
	if err != nil {
		return
	}

	log.Println(uid)

	s := reflect.ValueOf(ctx.Data["Errors"])

	if s.Len() == 0 {
		u := &User{Email: form.Email}
		db.First(u, u)
		if u.ID != 0 {
			sendPasswordResetEmail(u, "en-US")
			ctx.Data["Finished"] = true
		} else {
			ctx.Data["ShowSignUp"] = true

			ctx.Data["Errors"] = true
			ctx.Data["ErrorMsg"] = "We don't have this email in our database, please sign up again."
		}
	} else {
		ctx.Data["ErrorMsg"] = "Email field is required."
	}

	ctx.HTMLSet(200, "login", "passwordreset")
}

func passwordResetSignInView(ctx *macaron.Context, sess session.Store) {
	ctx.Data["Title"] = "Password Reset | "

	uid, err := decrypt([]byte(conf.DbPass[:16]), ctx.Params("uid"))
	if err != nil {
		return
	}

	u := &User{Email: uid}
	db.First(u, u)

	if u.ID == 0 {
		return
	}

	sess.Set("userID", u.ID)

	ctx.Redirect("/password-reset-finish/")
}

func passwordResetFinishView(ctx *macaron.Context) {
	ctx.Data["Title"] = "Password Reset Finish | "

	ctx.HTMLSet(200, "login", "passwordresetfinish")
}

func passwordResetFinishPostView(ctx *macaron.Context, form SignInForm) {
	ctx.Data["Title"] = "Password Reset Finish | "
	ctx.Data["Form"] = form
	ctx.Data["ShowOnlyPass"] = true

	s := reflect.ValueOf(ctx.Data["Errors"])

	if s.Len() == 0 {
		user := ctx.Data["User"].(*User)
		if user.Address == form.Address {
			hashedPassword, err := bcrypt.GenerateFromPassword([]byte(form.Password), 8)
			for err != nil {
				time.Sleep(200 * time.Microsecond)
				hashedPassword, err = bcrypt.GenerateFromPassword([]byte(form.Password), 8)
			}
			user.PasswordHash = string(hashedPassword)
			db.Save(user)

			ctx.Data["Finish"] = true
		} else {
			ctx.Data["ErrorMsg"] = "The address is wrong."
		}
	} else {
		ctx.Data["ErrorMsg"] = "New password field is required."
	}

	ctx.HTMLSet(200, "login", "passwordresetfinish")
}
