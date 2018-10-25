package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"io"
	"log"

	"github.com/go-macaron/session"
	"gopkg.in/macaron.v1"
)

func newPageData(ctx *macaron.Context, sess session.Store, f *session.Flash) {
	ctx.Data["URI"] = ctx.Req.RequestURI

	a := ctx.GetCookie("address")
	r := ctx.GetCookie("referral")

	ctx.Data["NodeAddress"] = conf.NodeAddress

	if len(a) > 0 {
		u := &User{Address: a}
		db.First(u, u)
		if u.ID == 0 {
			if len(r) > 0 {
				u.Referral = r
			}
			u.Nickname = a
			u.Email = a
			db.Create(u)
		}
		if u.ID != 0 && len(u.BitcoinAddr) == 0 {
			var err error
			u.BitcoinAddr, err = bg.getAddress()
			if err != nil {
				log.Printf("Error in bg.getAddress: %s", err)
			} else {
				db.Save(u)
			}
		}
		if u.ID != 0 && len(u.EtherAddr) == 0 {
			var err error
			u.EtherAddr, err = eg.getAddress(uint32(u.ID))
			if err != nil {
				log.Printf("Error in eg.getAddress: %s", err)
			}
			db.Save(u)
		}
		db.Model(u).Association("Badges").Find(&u.Badges)

		if u.Email == u.Address {
			u.Email = ""
		}

		if u.Nickname == u.Address {
			u.Nickname = ""
		}

		ctx.Data["User"] = u
		sess.Set("userID", u.ID)
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
