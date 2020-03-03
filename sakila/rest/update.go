package rest

import (
	"github.com/rohan5564/sakila-inventario/db"
	"github.com/rohan5564/sakila-inventario/sakila"
	MyErrors "github.com/rohan5564/sakila-inventario/sakila/errors"

	"io/ioutil"
	"net/http"

	"github.com/go-chi/render"
)

// UpdateFilm updates the film returned from the query search with the body request information
func UpdateFilm(w http.ResponseWriter, r *http.Request) {
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		render.Render(w, r, MyErrors.Err400)
		return
	}

	films, exists := r.Context().Value("films").(*sakila.QueryFilm)
	if !exists || len(films.Result) == 0 {
		render.Render(w, r, MyErrors.Err404)
		return
	}

	film := films.Result[0]

	if err := film.FromJSON(data); err != nil {
		render.Render(w, r, MyErrors.Err422)
		return
	}

	if err := db.UpdateFilm(film); err != nil {
		render.Render(w, r, MyErrors.Err404)
		return
	}

	w.WriteHeader(http.StatusAccepted)
}
