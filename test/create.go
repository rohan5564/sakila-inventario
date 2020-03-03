package test

import (
	. "github.com/rohan5564/sakila-inventario/test/data"

	"fmt"
	"testing"
)

// CreateFilms inserts multiple films into the database
func CreateFilms(t *testing.T) {
	for i, film := range FilmTests {
		film := film
		t.Run(fmt.Sprintf("%s", film.Title), func(t *testing.T) {
			t.Logf("test #%d: %s", i+1, film.Title)

			if _, ok, msg := testFilmURL("POST", film); !ok {
				t.Error(msg)
			}
		})
	}
}
