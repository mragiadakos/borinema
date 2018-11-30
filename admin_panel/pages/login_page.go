package pages

import (
	"github.com/gopherjs/vecty"
	h "github.com/gopherjs/vecty/elem"
	"github.com/gopherjs/vecty/event"
	"github.com/mragiadakos/borinema/admin_panel/actions"
	"github.com/mragiadakos/borinema/admin_panel/services"
	"github.com/mragiadakos/borinema/admin_panel/store"
)

type LoginPage struct {
	vecty.Core
	username string
	password string
	errStr   string
}

func (lp *LoginPage) onPassword(event *vecty.Event) {
	lp.password = event.Target.Get("value").String()
}

func (lp *LoginPage) onUsername(event *vecty.Event) {
	lp.username = event.Target.Get("value").String()
}

func (lp *LoginPage) onSubmit(event *vecty.Event) {
	go func() {
		authJson := services.AuthorizationJson{
			Username: lp.username,
			Password: lp.password,
		}
		println(authJson.Username, authJson.Password)
		as := services.AuthService{}
		token, errMsg := as.PostLogin(authJson)
		if errMsg != nil {
			lp.errStr = errMsg.Error
			vecty.Rerender(lp)
			return
		} else {
			lp.password = ""
			lp.username = ""
			lp.errStr = ""
		}
		as.SaveToken(token)
		store.Dispatch(&actions.SetIsAdmin{
			IsAdmin: true,
		})
		store.Dispatch(&actions.ToRedirect{
			ToRedirect: actions.PAGE_ADMIN_MAIN,
		})
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
