package db

import (
	"github.com/rohan5564/sakila-inventario/sakila"

	"errors"
)

// SearchAllFilms Executes a search query to get all films
func SearchAllFilms(result *sakila.FilmResult) error {
	query := `SELECT * FROM sakila_crud.film_data;`

	rows, err := sakila.Statement.Query(query)

	if err != nil {
		return err
	}
	defer rows.Close()

	if err := scanFilmResult(result, rows); err != nil {
		return err
	}

	err = rows.Err()
	if err != nil {
		return err
	}

	return nil
}

// SearchSomeFilms Executes a search query to get some films by applying a filter
func SearchSomeFilms(filter *sakila.Film, options sakila.FilmPagination, result *sakila.FilmResult, total *int16) error {
	query := `CALL sakila_crud.search_films(` +
		call(options.Limit,
			options.Offset,
			options.OrderBy,
			options.Order,
			filter.ID,
			filter.Description,
			filter.Year,
			filter.Lang,
			filter.OriginalLang,
			filter.RentalDuration,
			filter.RentalPrice,
			filter.Lenght,
			filter.ReplacementCost,
			filter.Rating) +
		`);
		CALL sakila_crud.search_films_count(` +
		call(filter.ID,
			filter.Description,
			filter.Year,
			filter.Lang,
			filter.OriginalLang,
			filter.RentalDuration,
			filter.RentalPrice,
			filter.Lenght,
			filter.ReplacementCost,
			filter.Rating) +
		`);`

	rows, err := sakila.Statement.Query(query)

	if err != nil {
		return err
	}
	defer rows.Close()

	if err := scanFilmResult(result, rows); err != nil {
		return err
	}

	if !rows.NextResultSet() {
		return errors.New("expected more result sets")
	}

	for rows.Next() {
		if err := rows.Scan(total); err != nil {
			return err
		}
	}

	return nil
}

// SearchFilm Executes a search query to get a specific film
func SearchFilm(filter *sakila.Film, result *sakila.FilmResult) error {
	query := `CALL sakila_crud.search_film(` + call(filter.ID, filter.Title) + `);`

	rows, err := sakila.Statement.Query(query)

	if err != nil {
		return err
	}
	defer rows.Close()

	if err := scanFilmResult(result, rows); err != nil {
		return err
	}

	err = rows.Err()
	if err != nil {
		return err
	}

	return nil
}

// SearchCategories Executes a search query to get all categories
func SearchCategories(result *sakila.CategoryResult) error {
	query := `CALL sakila_crud.search_category(NULL);`

	call(&query, sakila.Category{})

	rows, err := sakila.Statement.Query(query)

	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		tmp := new(sakila.Category)
		if err := rows.Scan(&tmp.ID, &tmp.Name, &tmp.Count); err != nil {
			return err
		}
		*result = append(*result, tmp)
	}

	err = rows.Err()
	if err != nil {
		return err
	}

	return nil
}

// SearchActors Executes a search query to get all actors
func SearchActors(result *sakila.ActorResult) error {
	query := `CALL sakila_crud.search_actor(NULL, NULL);`

	rows, err := sakila.Statement.Query(query)

	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		tmp := new(sakila.Actor)
		if err := rows.Scan(&tmp.ID, &tmp.FirstName, &tmp.LastName, &tmp.Films); err != nil {
			return err
		}
		*result = append(*result, tmp)
	}

	err = rows.Err()
	if err != nil {
		return err
	}

	return nil
}
