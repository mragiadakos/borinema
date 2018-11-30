package main

import (
	"github.com/mragiadakos/borinema/admin_panel/actions"
	"github.com/mragiadakos/borinema/admin_panel/services"
	"github.com/mragiadakos/borinema/admin_panel/store"
)

func SetMovies() {
	go func() {
		ms := services.MovieService{}
		pag := services.PaginationJson{}
		pag.Limit = -1
		mvs, _ := ms.GetMovies(pag)
		store.Dispatch(&actions.SetMovies{
			Movies: mvs,
		})
	}()
}
