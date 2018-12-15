package store

import (
	"github.com/mragiadakos/borinema/admin_panel/actions"
	"github.com/mragiadakos/borinema/admin_panel/services"
)

var (
	IsAdmin     = false
	CurrentPage = actions.PAGE_NOTHING
	Movies      = []services.MovieJson{}
)
var Listeners = NewListenerRegistry()

func init() {
	Register(onAction)
}
func findMoviesIndex(movies []services.MovieJson, id string) int {
	index := -1
	for i, v := range Movies {
		if id == v.ID {
			index = i
			break
		}
	}
	return index
}
func onAction(action interface{}) {
	println(action)
	switch a := action.(type) {
	case *actions.SetIsAdmin:
		IsAdmin = a.IsAdmin
	case *actions.ToRedirect:
		CurrentPage = a.ToRedirect
	case *actions.SetMovies:
		Movies = a.Movies
	case *actions.SetFirstMovieInList:
		Movies = append([]services.MovieJson{a.Movie}, Movies...)
	case *actions.RemoveMovieFromList:
		index := findMoviesIndex(Movies, a.MovieId)
		Movies = append(Movies[:index], Movies[index+1:]...)
	case *actions.AppendMoviesToList:
		Movies = append(Movies, a.Movies...)

	case *actions.SetMovieProgress:
		index := findMoviesIndex(Movies, a.ID)
		Movies[index].Progress = a.Progress
		Movies[index].State = a.State
		Movies[index].Filetype = a.Filetype
	default:
		return
	}
	Listeners.Fire()
}
