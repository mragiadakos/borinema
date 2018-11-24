package main

import (
	"github.com/gopherjs/vecty"
	h "github.com/gopherjs/vecty/elem"
	"github.com/gopherjs/vecty/event"
)

type LoginPage struct {
	vecty.Core
	username string
	password string
	errStr   string
}

func (lp *LoginPage) onPassword(event *vecty.Event) {
	lp.password = event.Target.Get("value").String()
	println(lp.password)
	vecty.Rerender(lp)
}

func (lp *LoginPage) onUsername(event *vecty.Event) {
	lp.username = event.Target.Get("value").String()
	println(lp.username)
	vecty.Rerender(lp)
}

func (lp *LoginPage) onSubmit(event *vecty.Event) {
	go func() {
		authJson := AuthorizationJson{
			Username: lp.username,
			Password: lp.password,
		}
		as := AuthService{}
		token, errMsg := as.postLogin(authJson)
		if errMsg != nil {
			lp.errStr = errMsg.Error
		} else {
			lp.password = ""
			lp.username = ""
			lp.errStr = ""
		}
		vecty.Rerender(lp)
		as.saveToken(token)
	}()
}

func (lp *LoginPage) Render() vecty.ComponentOrHTML {
	return h.Body(
		h.Div(
			h.Span(vecty.Text("Username:")),
			h.Input(vecty.Markup(
				event.Change(lp.onUsername),
			)),
			h.Break(),
			h.Span(vecty.Text("Password:")),
			h.Input(vecty.Markup(
				event.Change(lp.onPassword),
			)),
			h.Break(),
			vecty.If(len(lp.errStr) > 0, vecty.Text("Error: "+lp.errStr)),
			h.Break(),
			h.Button(
				vecty.Markup(
					event.Click(lp.onSubmit),
				),
				vecty.Text("Submit"),
			),
		),
	)
}
