package components

import (
	"github.com/gopherjs/vecty"
	h "github.com/gopherjs/vecty/elem"
	"github.com/gopherjs/vecty/event"
	"github.com/mragiadakos/borinema/admin_panel/actions"
	"github.com/mragiadakos/borinema/admin_panel/services"
	"github.com/mragiadakos/borinema/admin_panel/store"
)

type AddMovieComponent struct {
	vecty.Core
	link   string
	name   string
	errStr string
}

func (fmc *AddMovieComponent) onLink(event *vecty.Event) {
	fmc.link = event.Target.Get("value").String()
}

func (fmc *AddMovieComponent) onName(event *vecty.Event) {
	fmc.name = event.Target.Get("value").String()
}

func (fmc *AddMovieComponent) onSubmit(event *vecty.Event) {
	go func() {
		ms := services.MovieService{}
		amj := services.AddMovieJson{}
		amj.Name = fmc.name
		amj.Url = fmc.link
		movieId, errMsg := ms.AddMovie(amj)
		if errMsg != nil {
			fmc.errStr = errMsg.Error
			vecty.Rerender(fmc)
			return
		}
		go func() {
			mv, err := ms.GetMovie(movieId.ID)
			if err == nil {
				store.Dispatch(&actions.SetFirstMovieInList{
					Movie: *mv,
				})
			}
		}()
	}()
}

func (fmc *AddMovieComponent) Render() vecty.ComponentOrHTML {
	return h.Div(
		h.Heading4(vecty.Text("Upload movie")),
		h.Span(vecty.Text("Name:")),
		h.Input(vecty.Markup(
			event.Change(fmc.onName),
		)),
		h.Break(),
		h.Span(vecty.Text("Link:")),
		h.Input(vecty.Markup(
			event.Change(fmc.onLink),
		)),
		h.Break(),
		vecty.If(len(fmc.errStr) > 0, vecty.Text("Error: "+fmc.errStr)),
		h.Break(),
		h.Button(
			vecty.Markup(
				event.Click(fmc.onSubmit),
			),
			vecty.Text("Submit"),
		),
	)
}
