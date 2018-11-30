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
