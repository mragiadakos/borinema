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
				vecty.Markup(vecty.Class("columns")),
				h.Div(
					vecty.Markup(vecty.Class("column"), vecty.Class("is-two-thirds")),
					&components.MoviesComponent{Movies: store.Movies}),

				h.Div(
					vecty.Markup(vecty.Class("column")),
					&components.AddMovieComponent{},
					h.Break(),
					&components.PlayMovieComponent{},
				),
			),
		),
	)
}
