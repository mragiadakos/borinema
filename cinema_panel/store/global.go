package store

import (
	"github.com/mragiadakos/borinema/cinema_panel/actions"
)

var Listeners = NewListenerRegistry()

var (
	MovieID       = ""
	MovieName     = ""
	MovieFiletype = ""
)

func init() {
	Register(onAction)
}

func onAction(action interface{}) {
	println(action)
	switch a := action.(type) {
	case *actions.SetMovie:
		MovieID = a.MovieID
		MovieName = a.Name
		MovieFiletype = a.Filetype
	default:
		return
	}
	Listeners.Fire()
}
