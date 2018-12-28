package store

import (
	"github.com/mragiadakos/borinema/admin_panel/actions"
	"github.com/mragiadakos/borinema/admin_panel/services"
)

var (
	IsAdmin                           = false
	CurrentPage                       = actions.PAGE_NOTHING
	Movies                            = []services.MovieJson{}
	SelectedMovie *services.MovieJson = nil
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

func getSelectedMovie(movies []services.MovieJson) *services.MovieJson {
	for _, v := range movies {
		if v.Selected {
			return &v
		}
	}
	return nil
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
		SelectedMovie = getSelectedMovie(Movies)
	case *actions.SetFirstMovieInList:
		Movies = append([]services.MovieJson{a.Movie}, Movies...)
		SelectedMovie = getSelectedMovie(Movies)
	case *actions.RemoveMovieFromList:
		index := findMoviesIndex(Movies, a.MovieId)
		Movies = append(Movies[:index], Movies[index+1:]...)
		SelectedMovie = getSelectedMovie(Movies)
	case *actions.AppendMoviesToList:
		Movies = append(Movies, a.Movies...)
		SelectedMovie = getSelectedMovie(Movies)
	case *actions.SetMovieProgress:
		index := findMoviesIndex(Movies, a.ID)
		Movies[index].Progress = a.Progress
		Movies[index].State = a.State
		Movies[index].Filetype = a.Filetype
	case *actions.SelectMovieFromList:
		index := findMoviesIndex(Movies, a.ID)
		Movies[index].Selected = a.IsSelected
		if a.IsSelected {
			selectedMovie := Movies[index]
			SelectedMovie = &selectedMovie
		} else {
			SelectedMovie = nil
		}
		if a.IsSelected {
			for i, v := range Movies {
				if v.Selected && v.ID != a.ID {
					Movies[i].Selected = false
				}
			}
		}
	default:
		return
	}
	Listeners.Fire()
}
