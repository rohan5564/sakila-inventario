package test

import (
	"github.com/rohan5564/sakila-inventario/sakila"

	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"

	"github.com/go-chi/chi"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

// Router is a router used in testing, shouldn't be used in
// the main code
var Router *chi.Mux

func cmpString(s1, s2 string) bool {
	return strings.ToLower(s1) == strings.ToLower(s2)

}
func cmpNullString(s1, s2 sakila.NullString) bool {
	return (s1.Valid == s2.Valid) && (s1.Valid == false || cmpString(s1.String, s2.String))
}

func cmpNullInt32(n1, n2 sakila.NullInt32) bool {
	return (n1.Valid == n2.Valid) && (n1.Valid == false || n1.Int32 == n2.Int32)
}

func cmpNullTime(t1, t2 sakila.NullTime) bool {
	return (t1.Valid == t2.Valid) && (t1.Valid == false || t1.Time.Equal(t2.Time))
}

func filmOptions() []cmp.Option {
	return []cmp.Option{
		cmp.Comparer(cmpString),
		cmp.Comparer(cmpNullString),
		cmp.Comparer(cmpNullInt32),
		// cmp.Comparer(cmpNullTime),
		cmpopts.IgnoreFields(sakila.Film{}, "ID", "LastUpdate")}
}

func eq(ifc1, ifc2 interface{}) bool {
	t1 := reflect.TypeOf(ifc1).Name()
	t2 := reflect.TypeOf(ifc2).Name()
	if t1 != t2 {
		return false
	}

	switch ifc1.(type) {
	case sakila.Film:
		film1 := ifc1.(sakila.Film)
		film2 := ifc2.(sakila.Film)
		return cmp.Equal(film1, film2, filmOptions()...)
	case sakila.Category:
		cat1 := ifc1.(sakila.Category)
		cat2 := ifc2.(sakila.Category)
		return cmp.Equal(cat1, cat2)
	case sakila.Actor:
		actor1 := ifc1.(sakila.Actor)
		actor2 := ifc2.(sakila.Actor)
		return cmp.Equal(actor1, actor2)
	default:
		return cmp.Equal(ifc1, ifc2)
	}
}

func diff(ifc1, ifc2 interface{}) string {
	t1 := reflect.TypeOf(ifc1).Name()
	t2 := reflect.TypeOf(ifc2).Name()
	if t1 != t2 {
		return "Different value types"
	}

	switch ifc1.(type) {
	case sakila.Film:
		film1 := ifc1.(sakila.Film)
		film2 := ifc2.(sakila.Film)
		return cmp.Diff(film1, film2, filmOptions()...)
	case sakila.Category:
		cat1 := ifc1.(sakila.Category)
		cat2 := ifc2.(sakila.Category)
		return cmp.Diff(cat1, cat2)
	case sakila.Actor:
		actor1 := ifc1.(sakila.Actor)
		actor2 := ifc2.(sakila.Actor)
		return cmp.Diff(actor1, actor2)
	default:
		return cmp.Diff(ifc1, ifc2)
	}
}

func testFilmURL(method string, films ...sakila.Film) (result []byte, ok bool, message string) {
	if l := len(films); l > 2 {
		return nil, false, fmt.Sprintf("Max. number of parameters acepted: 2. Got %d", l)
	}

	var data []byte
	id := fmt.Sprintf("%d", films[0].ID)

	var url string
	var buffer *bytes.Buffer
	var statusCode int

	switch method {
	case "GET":
		url = "http://localhost/API/rest/inventory/"
		if films[0].ID == 0 {
			url += "search/q=" + strings.ReplaceAll(films[0].Title, " ", "+") +
				"/lim=1&orderby=id&ord=desc"
		} else {
			url += "film/" + id + "/" + strings.ReplaceAll(films[0].Title, " ", "_")
		}
		statusCode = http.StatusOK
	case "POST":
		url = "http://localhost/API/rest/inventory/"
		data, _ = json.Marshal(films[0])
		buffer = bytes.NewBuffer(data)
		statusCode = http.StatusCreated
	case "PUT":
		url = "http://localhost/API/rest/inventory/film/" +
			id + "/" + strings.ReplaceAll(films[0].Title, " ", "_")
		data, _ = json.Marshal(films[1])
		buffer = bytes.NewBuffer(data)
		statusCode = http.StatusAccepted
	case "DELETE":
		url = "http://localhost/API/rest/inventory/film/" +
			id + "/" + strings.ReplaceAll(films[0].Title, " ", "_")
		statusCode = http.StatusNoContent
	default:
		return nil, false, fmt.Sprintf("Can't load method \"%s\"", method)
	}

	var req *http.Request
	if buffer != nil {
		req, _ = http.NewRequest(method, url, buffer)
	} else {
		req, _ = http.NewRequest(method, url, nil)
	}
	rec := httptest.NewRecorder()
	Router.ServeHTTP(rec, req)
	if status := rec.Code; status != statusCode {
		return nil, false, fmt.Sprintf("HTTP status code %v", status)
	}

	return rec.Body.Bytes(), true, ""
}

func filteredFilm(data []byte) *sakila.Film {
	ifc := make(map[string]interface{})
	json.Unmarshal(data, &ifc)

	list, ok := ifc["data"].([]interface{})
	switch {
	case !ok || len(list) == 0:
		return nil
	case len(list) == 1:
		jsonString, _ := json.Marshal(list[0])
		return uniqueFilm(jsonString)
	}

	return nil
}

func uniqueFilm(data []byte) *sakila.Film {
	result := &sakila.Film{}
	if err := result.FromJSON(data); err != nil {
		return nil
	}
	return result
}
