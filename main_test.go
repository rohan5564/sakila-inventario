package main

import (
	"github.com/rohan5564/sakila-inventario/sakila"
	"github.com/rohan5564/sakila-inventario/test"

	"flag"
	"fmt"
	"github.com/go-chi/chi/middleware"
	"os"
	"testing"
)

type testFunc struct {
	Name     string
	Function func(*testing.T)
}

func TestMain(m *testing.M) {
	flag.Parse()
	dburl := fmt.Sprintf("%s:%s@%s(%s:%d)/%s?timeout=5s&parseTime=true&multiStatements=true",
		session.User, session.Pass, session.Protocol, session.IP, session.Port, session.Schema)
	StartDB(dburl)
	defer sakila.Statement.Close()

	test.Router = API(middleware.Recoverer, middleware.StripSlashes, middleware.NoCache)
	status := m.Run()
	os.Exit(status)
}

func TestAPI(t *testing.T) {
	tests := []testFunc{
		{"Ping", ping},
		{"Create film", test.CreateFilms},
		{"Get film", test.GetFilms},
		{"Update film", test.UpdateFilms},
		{"Delete film", test.DeleteFilms}}

	var i int
	for i = range tests {
		test := tests[i]
		if !t.Run(test.Name, test.Function) { // this by running tests linearly
			break
		}
	}

	i++
	if i != len(tests) {
		for _, test := range tests[i:] {
			t.Logf("Untested case: %s\n", test.Name)
		}
	}

}

func ping(t *testing.T) {
	if err := sakila.Statement.Ping(); err != nil {
		t.Error(err)
	}
}
