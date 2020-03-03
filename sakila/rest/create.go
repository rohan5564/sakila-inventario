package rest

import (
	"github.com/rohan5564/sakila-inventario/db"
	"github.com/rohan5564/sakila-inventario/sakila"
	MyErrors "github.com/rohan5564/sakila-inventario/sakila/errors"

	"io/ioutil"
	"net/http"

	"github.com/go-chi/render"
)

// CreateFilm creates a film from the body request and inserts it into the database
func CreateFilm(w http.ResponseWriter, r *http.Request) {
	film := &sakila.Film{}
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		render.Render(w, r, MyErrors.Err400)
		return
	}

	if err := film.FromJSON(data); err != nil {
		render.Render(w, r, MyErrors.Err422)
		return
	}

	if err := db.AddFilm(film); err != nil {
		render.Render(w, r, MyErrors.Err403)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

// CreateCategory creates a category from the body request and inserts it into the database
func CreateCategory(w http.ResponseWriter, r *http.Request) {
	category := &sakila.Category{}
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		render.Render(w, r, MyErrors.Err400)
		return
	}

	if err := category.FromJSON(data); err != nil {
		render.Render(w, r, MyErrors.Err422)
		return
	}

	if err := db.AddCategory(category); err != nil {
		render.Render(w, r, MyErrors.Err403)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

// CreateActor creates an actor from the body request and inserts it into the database
func CreateActor(w http.ResponseWriter, r *http.Request) {
	actor := &sakila.Actor{}
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		render.Render(w, r, MyErrors.Err400)
		return
	}

	if err := actor.FromJSON(data); err != nil {
		render.Render(w, r, MyErrors.Err422)
		return
	}

	if err := db.AddActor(actor); err != nil {
		render.Render(w, r, MyErrors.Err403)
		return
	}

	w.WriteHeader(http.StatusCreated)
}
