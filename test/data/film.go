package testData

import (
	"github.com/rohan5564/sakila-inventario/sakila"

	"encoding/json"
	"io/ioutil"
	"log"
	"os"
)

var FilmTests []sakila.Film

func init() {
	dir, _ := os.Getwd()

	if file, err := ioutil.ReadFile(dir + "/test/data/films.json"); err == nil {
		json.Unmarshal(file, &FilmTests)
	} else {
		log.Fatal(err)
	}
	// for _, v := range FilmTests {
	// 	log.Println(v)
	// }
}
