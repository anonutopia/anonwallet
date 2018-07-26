package main

import (
	"html/template"
	"strings"

	"github.com/go-macaron/cache"
	"github.com/go-macaron/i18n"
	"github.com/go-macaron/session"
	"gopkg.in/macaron.v1"

	_ "github.com/go-macaron/session/redis"
)

const (
	PROJECT_NAME = "Anonwallet"
)

func initMacaron() *macaron.Macaron {
	m := macaron.Classic()

	m.Use(i18n.I18n(i18n.Options{
		Langs: []string{"hr", "sr", "en-US"},
		Names: []string{"Hrvatski", "Srpski", "English"},
	}))

	ro := macaron.RenderOptions{
		Layout: "layout",
		Funcs: []template.FuncMap{map[string]interface{}{
			"obfuscate": func(args ...interface{}) template.HTML {
				email := args[0].(string)
				email = strings.Replace(email, "@", "<span style=\"display:none\">evilspam</span>@", 1)
				return template.HTML(email)

			},
		}},
	}

	m.Use(macaron.Renderers(ro, "login:templates/login"))

	m.Use(session.Sessioner(session.Options{
		Provider:       "redis",
		ProviderConfig: "addr=127.0.0.1:6379",
	}))

	m.Use(cache.Cacher())

	return m
}
