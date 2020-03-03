package db

import (
	"github.com/rohan5564/sakila-inventario/sakila"
)

// UpdateFilm Updates a film from the database
func UpdateFilm(film *sakila.Film) error {
	actors := sakila.NullJSONArray{Array: &film.Actors, Keys: []string{"firstname", "lastname"}}
	query := `CALL sakila_crud.update_film(` +
		call(film.ID,
			film.Title,
			film.Categories,
			actors,
			film.Description,
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
