package components

import (
	"github.com/gopherjs/vecty"
	h "github.com/gopherjs/vecty/elem"
	"github.com/gopherjs/vecty/event"
	"github.com/mragiadakos/borinema/admin_panel/actions"
	"github.com/mragiadakos/borinema/admin_panel/services"
	"github.com/mragiadakos/borinema/admin_panel/store"
)

type LogoutComponent struct {
	vecty.Core
}

func (lc *LogoutComponent) onLogout(event *vecty.Event) {
	go func() {
		as := services.AuthService{}
		as.RemoveToken()
		store.Dispatch(&actions.SetIsAdmin{
			IsAdmin: false,
		})
		store.Dispatch(&actions.ToRedirect{
			ToRedirect: actions.PAGE_ADMIN_LOGIN,
		})
	}()

}

func (lc *LogoutComponent) Render() vecty.ComponentOrHTML {
	return h.Div(
		h.Button(
			vecty.Markup(
				event.Click(lc.onLogout),
			),
			vecty.Text("Logout"),
		),
	)

}
