package test

import (
	. "github.com/rohan5564/sakila-inventario/test/data"

	"fmt"
	"testing"
)

// GetFilms sends multiple request url to get specific films as response
func GetFilms(t *testing.T) {
	for i := range FilmTests {
		film := &FilmTests[i]
		t.Run(fmt.Sprintf("%s", film.Title), func(t *testing.T) {
			t.Logf("test #%d: %s", i, film.Title)
			if data, ok, msg := testFilmURL("GET", *film); !ok {
				t.Error(msg)
			} else {
				result := filteredFilm(data)
				switch {
				case result == nil:
					fmt.Println(string(data))
					t.Error("Film poiter is Nil, check the number of values in the body response")
				case !eq(*film, *result):
					t.Error(diff(*film, *result))
				default:
					*film = *result
				}
			}
		})
	}
}
