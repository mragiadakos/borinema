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

func onAction(action interface{}) {
	println(action)
	switch a := action.(type) {
	case *actions.SetIsAdmin:
		IsAdmin = a.IsAdmin
	case *actions.ToRedirect:
		CurrentPage = a.ToRedirect
	case *actions.SetMovies:
		Movies = a.Movies
	default:
		return
	}
	Listeners.Fire()
}
