package test

import (
	. "github.com/rohan5564/sakila-inventario/test/data"

	"fmt"
	"testing"
)

// DeleteFilms removes multiple films from the database
func DeleteFilms(t *testing.T) {
	for i, film := range FilmTests {
		film := film // to support parallel task
		t.Run(fmt.Sprintf("%s", film.Title), func(t *testing.T) {
			t.Logf("test #%d: %s", i, film.Title)
			if _, ok, msg := testFilmURL("DELETE", film); !ok {
				t.Error(msg)
			}
		})
	}
}
