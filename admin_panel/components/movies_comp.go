package components

import (
	"github.com/gopherjs/vecty"
	h "github.com/gopherjs/vecty/elem"
	"github.com/mragiadakos/borinema/admin_panel/services"
)

type MoviesComponent struct {
	vecty.Core
	Movies []services.MovieJson
}

func (mc *MoviesComponent) Render() vecty.ComponentOrHTML {
	var lis vecty.List
	for _, v := range mc.Movies {
		item := h.ListItem(
			h.Span(vecty.Text(v.Name)),
		)
		lis = append(lis, item)
	}
	println("render")
	return h.Div(
		h.Heading4(vecty.Text("Movies")),
		h.UnorderedList(lis),
	)

}
