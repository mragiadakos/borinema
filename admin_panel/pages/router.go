package pages

import (
	"github.com/gopherjs/vecty"
	h "github.com/gopherjs/vecty/elem"
	"github.com/mragiadakos/borinema/admin_panel/actions"
	"marwan.io/vecty-router"
)

type Router struct {
	vecty.Core
}

func (r *Router) Render() vecty.ComponentOrHTML {
	return h.Body(
		router.NewRoute(actions.PAGE_ADMIN_LOGIN.String(), &LoginPage{}, router.NewRouteOpts{ExactMatch: true}),
		router.NewRoute(actions.PAGE_ADMIN_MAIN.String(), &AdminPage{}, router.NewRouteOpts{ExactMatch: true}),
	)
}
