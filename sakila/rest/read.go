package rest

import (
	"github.com/rohan5564/sakila-inventario/sakila"
	MyErrors "github.com/rohan5564/sakila-inventario/sakila/errors"

	"log"
	"net/http"
	"runtime"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
)

// listrenderer is a wrapper that contains the information related to
// films and pagination
type listRenderer struct {
	Data        sakila.FilmResult        `json:"data"`
	Information sakila.ResultInformation `json:"information"`
}

// Render implements the renderer interface from go-chi/render
func (lr listRenderer) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

// GetFilm prints the films returned from the query search into the response data
func GetFilm(w http.ResponseWriter, r *http.Request) {
	// w.Header().Set("Access-Control-Allow-Origin", "*")
	films, assert := r.Context().Value("films").(*sakila.QueryFilm)
	if !assert || films == nil {
		render.Render(w, r, MyErrors.Err404)
		return
	}

	// github.com/go-chi/chi/blob/master/middleware/logger.go#L66
	logger := middleware.GetLogEntry(r) != nil

	films.GetResult()
	switch {
	case logger: // just for testing
		var mem runtime.MemStats
		runtime.ReadMemStats(&mem)
		log.Println(mem.Alloc/1024, " kb with ", len(films.Result), " films found")
		fallthrough
	case len(films.Result) > 0:
		if films.Info.SingleTitle {
			render.Render(w, r, films.Result[0])
		} else {
			render.Render(w, r, listRenderer{films.Result, films.Info})
		}
		w.WriteHeader(http.StatusOK)
	default:
		render.Render(w, r, MyErrors.Err404)
	}
}

// GetCategories prints the list of categories avairables
func GetCategories(w http.ResponseWriter, r *http.Request) {
	// w.Header().Set("Access-Control-Allow-Origin", "*")
	categories, assert := r.Context().Value("categories").(*sakila.CategoryResult)
	if !assert || len(*categories) == 0 {
		render.Render(w, r, MyErrors.Err404)
		return
	}

	list := []render.Renderer{}
	for _, category := range *categories {
		list = append(list, category)
	}

	render.RenderList(w, r, list)
}

// GetActors prints the list of categories avairables
func GetActors(w http.ResponseWriter, r *http.Request) {
	// w.Header().Set("Access-Control-Allow-Origin", "*")
	actors, assert := r.Context().Value("actors").(*sakila.ActorResult)
	if !assert || len(*actors) == 0 {
		render.Render(w, r, MyErrors.Err404)
		return
	}

	list := []render.Renderer{}
	for _, actor := range *actors {
		list = append(list, actor)
	}

	render.RenderList(w, r, list)
}
