package test

import (
	"github.com/rohan5564/sakila-inventario/sakila"
	. "github.com/rohan5564/sakila-inventario/test/data"

	"fmt"
	"testing"
)

// UpdateFilms updates the values from multiple films into the database
func UpdateFilms(t *testing.T) {
	for i := range FilmTests {
		film := &FilmTests[i] // to use referenced values and support parallel task
		t.Run(fmt.Sprintf("%s", film.Title), func(t *testing.T) {
			t.Logf("test #%d: %s", i, film.Title)
			old := *film
			film.Description = &sakila.NullString{String: "modified field", Valid: true}

			if _, ok, msg := testFilmURL("PUT", old, *film); !ok {
				t.Error(msg)
			} else if data, ok, msg2 := testFilmURL("GET", *film); ok {
				result := uniqueFilm(data)
				switch {
				case result == nil:
					t.Error("Film poiter is Nil, check the number of values in the body response")
				case eq(old, *result):
					t.Errorf("Film %s not updated", old.Title)
				}
			} else {
				t.Error(msg2)
			}
		})
	}
}
