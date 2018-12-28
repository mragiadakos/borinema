package main

import (
	"github.com/gopherjs/vecty"
	"github.com/mragiadakos/borinema/admin_panel/actions"
	"github.com/mragiadakos/borinema/admin_panel/pages"
	"github.com/mragiadakos/borinema/admin_panel/services"
	"github.com/mragiadakos/borinema/admin_panel/store"
	router "marwan.io/vecty-router"
)

func main() {
	r := &pages.Router{}
	vecty.SetTitle("Borinema's Admin Page")
	vecty.AddStylesheet("/admin_panel/node_modules/bulma/css/bulma.min.css")
	vecty.AddStylesheet("/admin_panel/admin_panel.css")
	vecty.RenderBody(r)
	EnableWebsocket()
	store.Listeners.Add(r, func() {
		if len(store.CurrentPage) > 0 {
			router.Redirect(store.CurrentPage.String())

			if store.CurrentPage == actions.PAGE_ADMIN_MAIN {
				SetMovies()
			}
			store.Dispatch(&actions.ToRedirect{
				ToRedirect: actions.PAGE_NOTHING,
			})
			return
		}

		vecty.RenderBody(r)
	})

	go func() {
		as := services.AuthService{}
		isAdmin, _ := as.GetIsAdmin()
		store.Dispatch(&actions.SetIsAdmin{
			IsAdmin: isAdmin,
		})
		if isAdmin {
			store.Dispatch(&actions.ToRedirect{
				ToRedirect: actions.PAGE_ADMIN_MAIN,
			})

		} else {
			store.Dispatch(&actions.ToRedirect{
				ToRedirect: actions.PAGE_ADMIN_LOGIN,
			})
		}
	}()
}
