package actions

import "github.com/mragiadakos/borinema/admin_panel/services"

type SetIsAdmin struct {
	IsAdmin bool
}

type ToRedirect struct {
	ToRedirect PageType
}

type SetMovies struct {
	Movies []services.MovieJson
}

type SetMovieProgress struct {
	ID       string
	Progress float64
	State    string
	Filetype string
}

type SelectMovieFromList struct {
	ID         string
	IsSelected bool
}
type SetFirstMovieInList struct {
	Movie services.MovieJson
}

type RemoveMovieFromList struct {
	MovieId string
}

type AppendMoviesToList struct {
	Movies []services.MovieJson
}
