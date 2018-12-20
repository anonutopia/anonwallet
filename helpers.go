package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/go-macaron/session"
	macaron "gopkg.in/macaron.v1"
)

func newPageData(ctx *macaron.Context, sess session.Store, f *session.Flash) {
	if !strings.HasSuffix(ctx.Req.URL.Path, "/") {
		http.Redirect(ctx.Resp, ctx.Req.Request, ctx.Req.URL.Path+"/", http.StatusFound)
		return
	}

	ctx.Data["URI"] = ctx.Req.RequestURI
	ctx.Data["NodeAddress"] = conf.NodeAddress
	ctx.Data["Notification"] = false
	ctx.Data["NotificationTitle"] = ""
	ctx.Data["NotificationMessage"] = ""

	userID := sess.Get("userID")
	if userID != nil {
		user := &User{}
		db.First(user, userID)
		ctx.Data["User"] = user

		if !user.EmailVerified {
			ctx.Data["Notification"] = true
			ctx.Data["NotificationTitle"] = "Please Verify Your Email Address"
			ctx.Data["NotificationMessage"] = "Please verify your email address and you will receive free anotes. If you haven't received our email, please check your spam folder."
		} else if !user.ReceivedFreeAnote && user.FacebookShareEnabled() {
			ctx.Data["Notification"] = true
			ctx.Data["NotificationTitle"] = "Get More Free Anotes"
			ctx.Data["NotificationMessage"] = "Get more free anotes by inviting your friends over Facebook (use settings page to share or profit page to get referral URL) or join our Telegram group and ask our bot for it."
		} else if !user.ReceivedFreeAnote {
			ctx.Data["Notification"] = true
			ctx.Data["NotificationTitle"] = "Get More Free Anotes"
			ctx.Data["NotificationMessage"] = "Get more free anotes by joining our Telegram group and asking our bot for it."
		} else if user.FacebookShareEnabled() {
			ctx.Data["Notification"] = true
			ctx.Data["NotificationTitle"] = "Get More Free Anotes"
			ctx.Data["NotificationMessage"] = "Get more free anotes by inviting your friends over Facebook (use settings page or to share or profit page to get referral URL)."
		}

		if user.ID != 0 && len(user.BitcoinAddr) == 0 {
			var err error
			user.BitcoinAddr, err = bg.getAddress()
			if err != nil {
				log.Printf("Error in bg.getAddress: %s", err)
				logTelegram(fmt.Sprintf("Error in bg.getAddress: %s", err))
			} else {
				db.Save(user)
			}
		}
		if user.ID != 0 && len(user.EtherAddr) == 0 {
			var err error
			user.EtherAddr, err = eg.getAddress(uint32(user.ID))
			if err != nil {
				log.Printf("Error in eg.getAddress: %s", err)
				logTelegram(fmt.Sprintf("Error in eg.getAddress: %s", err))
			}
			db.Save(user)
		}
		db.Model(user).Association("Badges").Find(&user.Badges)
	}

	ctx.Data["Anote"] = anote
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

func encrypt(key []byte, message string) (encmess string, err error) {
	plainText := []byte(message)

	block, err := aes.NewCipher(key)
	if err != nil {
		return
	}

	//IV needs to be unique, but doesn't have to be secure.
	//It's common to put it at the beginning of the ciphertext.
	cipherText := make([]byte, aes.BlockSize+len(plainText))
	iv := cipherText[:aes.BlockSize]
	if _, err = io.ReadFull(rand.Reader, iv); err != nil {
		return
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(cipherText[aes.BlockSize:], plainText)

	//returns to base64 encoded string
	encmess = base64.URLEncoding.EncodeToString(cipherText)
	return
}

func decrypt(key []byte, securemess string) (decodedmess string, err error) {
	cipherText, err := base64.URLEncoding.DecodeString(securemess)
	if err != nil {
		return
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return
	}

	if len(cipherText) < aes.BlockSize {
		err = errors.New("Ciphertext block size is too short!")
		return
	}

	//IV needs to be unique, but doesn't have to be secure.
	//It's common to put it at the beginning of the ciphertext.
	iv := cipherText[:aes.BlockSize]
	cipherText = cipherText[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)
	// XORKeyStream can work in-place if the two arguments are the same.
	stream.XORKeyStream(cipherText, cipherText)

	decodedmess = string(cipherText)
	return
}
