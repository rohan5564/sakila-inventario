package db

import (
	"github.com/rohan5564/sakila-inventario/sakila"

	"strings"
)

// AddFilm Inserts a film into the database
func AddFilm(film *sakila.Film) error {
	actors := sakila.NullJSONArray{Array: &film.Actors, Keys: []string{"firstname", "lastname"}}
	query := `CALL sakila_crud.add_film(` +
		call(film.Title,
			film.Description,
			film.Categories,
			actors,
			film.Year,
			film.Lang,
			film.OriginalLang,
			film.RentalDuration,
			film.RentalPrice,
			film.Lenght,
			film.ReplacementCost,
			film.Rating,
			film.SpecialFeatures) +
		`);`

	if _, err := sakila.Statement.Exec(query); err != nil {
		return err
	}

	return nil
}

// AddFilmCategory Adds categories to a film
func AddFilmCategory(film *sakila.Film, categories []string) error {
	for _, category := range categories {
		query := `CALL sakila_crud.add_film_category(` +
			call(film.Title) +
			`'` + category + `');`
		call(&query, *film)
		if _, err := sakila.Statement.Exec(string(query)); err != nil {
			return err
		}
	}
	return nil

}

// AddCategory Inserts a category into the database
func AddCategory(category *sakila.Category) error {
	query := `CALL sakila_crud.add_category('` + category.Name + `');`
	if _, err := sakila.Statement.Exec(query); err != nil {
		return err
	}
	return nil
}

// AddFilmActor Adds actors to a film
func AddFilmActor(film *sakila.Film, actors []string) error {
	for _, actor := range actors {
		arr := strings.Split(actor, " ")
		firstName := arr[0]
		lastName := arr[1]
		query := `CALL sakila_crud.add_film_actor(` +
			call(film.Title) +
			`,'` + firstName + `','` +
			lastName + `');`
		if _, err := sakila.Statement.Exec(query); err != nil {
			return err
		}
	}
	return nil

}

// AddActor Inserts an actor into the database
func AddActor(actor *sakila.Actor) error {
	query := `CALL sakila_crud.add_actor(
		'` + actor.FirstName + `','` +
		actor.LastName + `');`
	if _, err := sakila.Statement.Exec(query); err != nil {
		return err
	}
	return nil
}
