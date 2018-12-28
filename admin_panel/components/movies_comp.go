package components

import (
	"strconv"

	"github.com/mragiadakos/borinema/admin_panel/actions"
	"github.com/mragiadakos/borinema/admin_panel/store"

	"github.com/gopherjs/vecty"
	h "github.com/gopherjs/vecty/elem"
	"github.com/gopherjs/vecty/event"
	"github.com/gopherjs/vecty/prop"
	"github.com/mragiadakos/borinema/admin_panel/services"
)

type MoviesComponent struct {
	vecty.Core
	Movies []services.MovieJson `vecty:"prop"`
}

func (mc *MoviesComponent) renderMoviesTableHead() vecty.ComponentOrHTML {
	return h.TableHead(
		h.TableRow(
			h.TableHeader(h.Abbreviation(vecty.Markup(vecty.Attribute("title", "Name")), vecty.Text("Name"))),
			h.TableHeader(h.Abbreviation(vecty.Markup(vecty.Attribute("title", "File Type")), vecty.Text("File Type"))),
			h.TableHeader(h.Abbreviation(vecty.Markup(vecty.Attribute("title", "State")), vecty.Text("State"))),
			h.TableHeader(h.Abbreviation(vecty.Markup(vecty.Attribute("title", "Progress")), vecty.Text("Progress"))),
			h.TableHeader(h.Abbreviation(vecty.Markup(vecty.Attribute("title", "Delete")), vecty.Text("Delete"))),
			h.TableHeader(h.Abbreviation(vecty.Markup(vecty.Attribute("title", "Play")), vecty.Text("Play"))),
		),
	)
}

func (mc *MoviesComponent) renderMoviesTableRow(index int, movie services.MovieJson) vecty.ComponentOrHTML {
	prog := strconv.FormatFloat(movie.Progress, 'f', 1, 64)
	return h.TableRow(
		h.TableData(vecty.Text(movie.Name)),
		h.TableData(vecty.Text(movie.Filetype)),
		h.TableData(vecty.Text(movie.State)),
		h.TableData(vecty.Text(prog)),
		h.TableData(h.Button(
			vecty.Markup(
				event.Click(mc.onDeleteMovie(movie.ID)),
			),
			vecty.Text("x"))),
		h.TableData(h.Input(
			vecty.Markup(
				prop.Type(prop.TypeCheckbox),
				prop.Checked(movie.Selected),
				event.Change(mc.onSelectingMovie(movie.ID, movie.Selected)),
			),
		)),
	)
}

func (mc *MoviesComponent) onDeleteMovie(id string) func(e *vecty.Event) {
	return func(e *vecty.Event) {
		go func() {
			ms := services.MovieService{}
			ms.DeleteMovie(id)
			store.Dispatch(&actions.RemoveMovieFromList{MovieId: id})
		}()
	}
}

func (mc *MoviesComponent) onSelectingMovie(id string, selected bool) func(e *vecty.Event) {
	return func(e *vecty.Event) {
		go func() {
			ms := services.MovieService{}
			if selected {
				err := ms.RemoveMovieSelection()
				if err != nil {
					println("Error: " + err.Error)
				}
			} else {
				err := ms.SelectMovie(id)
				if err != nil {
					println("Error: " + err.Error)
				}
			}
			store.Dispatch(&actions.SelectMovieFromList{ID: id, IsSelected: !selected})
		}()
	}
}

func (mc *MoviesComponent) onMore(e *vecty.Event) {
	go func() {
		if len(store.Movies) > 0 {
			createdAt := store.Movies[len(store.Movies)-1].CreatedAt
			ms := services.MovieService{}
			pag := services.PaginationJson{}
			pag.Limit = 10
			pag.LastSeenDate = &createdAt
			mvs, _ := ms.GetMovies(pag)
			store.Dispatch(&actions.AppendMoviesToList{Movies: mvs})
		}
	}()
}

func (mc *MoviesComponent) Render() vecty.ComponentOrHTML {
	var lis vecty.List
	for i, v := range mc.Movies {
		item := mc.renderMoviesTableRow(i, v)
		lis = append(lis, item)
	}
	println("render")
	return h.Div(
		h.Heading1(vecty.Text("Movies")),
		h.Div(
			vecty.Markup(vecty.Class("scrollbar")),
			h.Table(
				vecty.Markup(vecty.Class("table"), vecty.Class("is-fullwidth")),
				mc.renderMoviesTableHead(),
				lis,
			),
			h.Button(vecty.Markup(event.Click(mc.onMore)),
				vecty.Text("More")),
		),
	)

}
