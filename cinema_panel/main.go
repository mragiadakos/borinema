package main

import (
	"github.com/gopherjs/vecty"
	"github.com/mragiadakos/borinema/cinema_panel/store"
)

func main() {
	cp := &CinemaPage{}
	vecty.SetTitle("Borinema's Cinema Page")
	vecty.AddStylesheet("/admin_panel/node_modules/bulma/css/bulma.min.css")
	vecty.RenderBody(cp)
	enableWebsocket()
	getMovieOnStart()
	store.Listeners.Add(cp, func() {
		vecty.RenderBody(cp)
	})
}
