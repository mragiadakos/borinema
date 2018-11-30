package components

import (
	"github.com/gopherjs/vecty"
	h "github.com/gopherjs/vecty/elem"
	"github.com/gopherjs/vecty/event"
	"github.com/mragiadakos/borinema/admin_panel/actions"
	"github.com/mragiadakos/borinema/admin_panel/services"
	"github.com/mragiadakos/borinema/admin_panel/store"
)

type FormMovieComponent struct {
	vecty.Core
	link   string
	name   string
	errStr string
}

func (fmc *FormMovieComponent) onLink(event *vecty.Event) {
	fmc.link = event.Target.Get("value").String()
}

func (fmc *FormMovieComponent) onName(event *vecty.Event) {
	fmc.name = event.Target.Get("value").String()
}

func (fmc *FormMovieComponent) onSubmit(event *vecty.Event) {
	go func() {
		ms := services.MovieService{}
		amj := services.AddMovieJson{}
		amj.Name = fmc.name
		amj.Url = fmc.link
		_, errMsg := ms.AddMovie(amj)
		if errMsg != nil {
			fmc.errStr = errMsg.Error
			vecty.Rerender(fmc)
			return
		}
		go func() {
			pag := services.PaginationJson{}
			pag.Limit = -1
			mvs, _ := ms.GetMovies(pag)
			store.Dispatch(&actions.SetMovies{
				Movies: mvs,
			})
		}()
	}()
}

func (fmc *FormMovieComponent) Render() vecty.ComponentOrHTML {
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
