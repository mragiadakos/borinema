package main

import (
	"github.com/gopherjs/vecty"
	h "github.com/gopherjs/vecty/elem"
	"github.com/mragiadakos/borinema/cinema_panel/store"
)

type CinemaPage struct {
	vecty.Core
}

func (cp *CinemaPage) Render() vecty.ComponentOrHTML {
	hasMovie := len(store.MovieID) > 0
	return h.Body(
		h.Div(
			h.Heading4(vecty.Text("Web Cinema")),
			vecty.If(hasMovie,
				h.Div(
					h.Span(vecty.Text("Title: "+store.MovieName)),
					h.Break(),
					h.Video(
						vecty.Markup(
							vecty.Attribute("oncontextmenu", "return false;"),
							vecty.Property("width", 400),
							vecty.Property("id", "media-video"),
						),
						h.Source(
							vecty.Markup(
								vecty.Property("src", "/api/cinema/movie/"+store.MovieID),
								vecty.Property("type", "video/"+store.MovieFiletype),
							),
						),
					),
					h.Break(),
				),
			),
			vecty.If(!hasMovie,
				h.Div(
					h.Span(vecty.Text("No movie have been selected.")),
					h.Break(),
				),
			),
		))
}
