package rest

import (
	"github.com/rohan5564/sakila-inventario/db"
	"github.com/rohan5564/sakila-inventario/sakila"
	MyErrors "github.com/rohan5564/sakila-inventario/sakila/errors"

	"net/http"

	"github.com/go-chi/render"
)

// DeleteFilm Removes a film inside the request context from the database
func DeleteFilm(w http.ResponseWriter, r *http.Request) {
	films, ok := r.Context().Value("films").(*sakila.QueryFilm)
	if !ok || films == nil || len(films.Result) == 0 {
		render.Render(w, r, MyErrors.Err404)
		return
	}
	film := films.Result[0]
	if err := db.RemoveFilm(film); err != nil {
		render.Render(w, r, MyErrors.Err404)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// DeleteCategory Removes a category inside the request context from the database
func DeleteCategory(w http.ResponseWriter, r *http.Request) {
	category, ok := r.Context().Value("category").(*sakila.Category)
	if !ok || category == nil {
		render.Render(w, r, MyErrors.Err404)
		return
	}

	if err := db.RemoveCategory(category); err != nil {
		render.Render(w, r, MyErrors.Err404)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// DeleteActor Removes an actor inside the request context from the database
func DeleteActor(w http.ResponseWriter, r *http.Request) {
	actor, ok := r.Context().Value("actor").(*sakila.Actor)
	if !ok || actor == nil {
		render.Render(w, r, MyErrors.Err404)
		return
	}

	if err := db.RemoveActor(actor); err != nil {
		render.Render(w, r, MyErrors.Err404)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
