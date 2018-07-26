package main

import macaron "gopkg.in/macaron.v1"

type Locale struct {
	Strings map[string]string
}

func initLocale(ctx *macaron.Context) map[string]string {
	// l := &Locale{}
	// l.Strings = make(map[string]string)
	l := make(map[string]string)

	// l.Strings["jsFieldNotEmpty"] = ctx.Tr("jsFieldNotEmpty")
	l["jsFieldNotEmpty"] = ctx.Tr("jsFieldNotEmpty")

	return l
}
