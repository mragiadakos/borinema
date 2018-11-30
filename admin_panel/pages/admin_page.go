package pages

import (
	"github.com/gopherjs/vecty"
	h "github.com/gopherjs/vecty/elem"
	"github.com/mragiadakos/borinema/admin_panel/components"
	"github.com/mragiadakos/borinema/admin_panel/store"
)

type AdminPage struct {
	vecty.Core
}

func (mp *AdminPage) Render() vecty.ComponentOrHTML {
	return h.Body(
		h.Div(
			vecty.If(store.IsAdmin,
				&components.LogoutComponent{},
			),
			h.Break(),
			h.Div(
				&components.FormMovieComponent{},
			),
			h.Break(),
			h.Div(&components.MoviesComponent{Movies: store.Movies}),
		),
	)
}
