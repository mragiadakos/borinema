package components

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/gopherjs/vecty"
	h "github.com/gopherjs/vecty/elem"
	"github.com/gopherjs/vecty/event"
	"github.com/mragiadakos/borinema/admin_panel/store"
)

type PlayMovieComponent struct {
	vecty.Core
}

func (pmc *PlayMovieComponent) onPlay(event *vecty.Event) {
	video := js.Global.Get("document").Call("getElementById", "media-video")
	video.Call("play")
}

func (pmc *PlayMovieComponent) onPause(event *vecty.Event) {
	video := js.Global.Get("document").Call("getElementById", "media-video")
	video.Call("pause")
}

func (pmc *PlayMovieComponent) onStop(event *vecty.Event) {
	video := js.Global.Get("document").Call("getElementById", "media-video")
	video.Call("pause")
	video.Set("currentTime", 0)
}

func (pmc *PlayMovieComponent) Render() vecty.ComponentOrHTML {
	movieIsSelected := store.SelectedMovie != nil
	movieSelectedTitle := ""
	filetype := ""
	movieId := ""
	if movieIsSelected {
		movieSelectedTitle = store.SelectedMovie.Name
		filetype = store.SelectedMovie.Filetype
		movieId = store.SelectedMovie.ID
	}
	return h.Div(
		h.Heading4(vecty.Text("Play movie")),
		vecty.If(movieIsSelected,
			h.Div(
				h.Span(vecty.Text("Title: "+movieSelectedTitle)),
				h.Break(),
				h.Video(
					vecty.Markup(
						vecty.Attribute("oncontextmenu", "return false;"),
						vecty.Property("width", 400),
						vecty.Property("id", "media-video"),
					),
					h.Source(
						vecty.Markup(
							vecty.Property("src", "/api/cinema/movie/"+movieId),
							vecty.Property("type", "video/"+filetype),
						),
					),
				),
				h.Break(),
				h.Button(
					vecty.Markup(
						event.Click(pmc.onPlay),
					),
					vecty.Text("play"),
				),
				h.Button(
					vecty.Markup(
						event.Click(pmc.onStop),
					),
					vecty.Text("stop"),
				),
				h.Button(
					vecty.Markup(
						event.Click(pmc.onPause),
					),
					vecty.Text("pause"),
				),
			),
		),
		vecty.If(!movieIsSelected,
			h.Div(
				h.Span(vecty.Text("No movie have been selected.")),
				h.Break(),
			),
		),
	)
}
