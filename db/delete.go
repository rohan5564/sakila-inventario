package db

import (
	"github.com/rohan5564/sakila-inventario/sakila"
)

// RemoveFilm Removes a film from the database
func RemoveFilm(film *sakila.Film) error {
	query := `CALL sakila_crud.remove_film(` + call(film.ID) + `);`

	if _, err := sakila.Statement.Exec(query); err != nil {
		return err
	}

	return nil
}

// RemoveFilmCategory Removes a list of categories from a film
func RemoveFilmCategory(film *sakila.Film, categories []string) error {
	for _, category := range categories {
		query := `CALL sakila_crud.remove_category_film(` +
			call(film.ID) +
			`'` + category + `');`
		if _, err := sakila.Statement.Exec(query); err != nil {
			return err
		}
	}
	return nil
}

// RemoveFilmActor Removes a list of actors from a film
func RemoveFilmActor(film *sakila.Film, actors []string) error {
	for _, actor := range actors {
		query := `CALL sakila_crud.remove_actor_film(` +
			call(film.ID) +
			`'` + actor + `');`
		if _, err := sakila.Statement.Exec(query); err != nil {
			return err
		}
	}
	return nil
}

// RemoveCategory Removes a category from the database
func RemoveCategory(category *sakila.Category) error {
	query := `CALL sakila_crud.remove_category(` + call(category.ID) + `);`
	if _, err := sakila.Statement.Exec(query); err != nil {
		return err
	}
	return nil
}

// RemoveActor Removes an actor from the database
func RemoveActor(actor *sakila.Actor) error {
	query := `CALL sakila_crud.remove_actor(` + call(actor.ID) + `);`
	if _, err := sakila.Statement.Exec(query); err != nil {
		return err
	}
	return nil
}
