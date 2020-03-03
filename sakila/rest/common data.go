package rest

import (
	"github.com/rohan5564/sakila-inventario/db"
	"github.com/rohan5564/sakila-inventario/sakila"
	MyErrors "github.com/rohan5564/sakila-inventario/sakila/errors"

	"context"
	"database/sql"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
)

const (
	PathTo_Inventory  = "/API/rest/inventory"
	PathTo_Categories = "/API/rest/categories"
	PathTo_Actors     = "/API/rest/actors"
	id                = "id"
	title             = "title"
	search            = "q"
	year              = "year"
	language          = "lang"
	rentalDuration    = "rdur"
	price             = "price"
	duration          = "len"
	replacementCost   = "replcos"
	rating            = "rating"
	specialFeatures   = "spec"
	lastUpdate        = "lu"
	limit             = "lim"
	offset            = "page"
	orderby           = "orderby"
	order             = "ord"
)

// NewFilm Middleware that checks if the connection to database is alive
// and adds an initial film variable into a context value
func NewFilm(sig http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := sakila.Statement.Ping(); err != nil {
			render.Render(w, r, MyErrors.Err500)
			return
		}
		options := sakila.FilmPagination{Limit: 20}
		information := sakila.ResultInformation{}

		films := &sakila.QueryFilm{
			Options: options,
			Info:    information}

		ctx := context.WithValue(r.Context(), "films", films)
		sig.ServeHTTP(w, r.WithContext(ctx))
	})
}

// NewCategory Middleware that checks if the connection to database is alive
// and adds an initial category variable into a context value
func NewCategory(sig http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := sakila.Statement.Ping(); err != nil {
			render.Render(w, r, MyErrors.Err500)
			return
		}
		categories := &sakila.CategoryResult{}
		ctx := context.WithValue(r.Context(), "categories", categories)
		sig.ServeHTTP(w, r.WithContext(ctx))
	})
}

// NewActor Middleware that checks if the connection to database is alive
// and adds an initial actor variable into a context value
func NewActor(sig http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := sakila.Statement.Ping(); err != nil {
			render.Render(w, r, MyErrors.Err500)
			return
		}
		actors := &sakila.ActorResult{}
		ctx := context.WithValue(r.Context(), "actors", actors)
		sig.ServeHTTP(w, r.WithContext(ctx))
	})
}

// FilmRequest Middleware that runs the inventory search to find one or
// many films (with or without filter)
func FilmRequest(sig http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		films, _ := r.Context().Value("films").(*sakila.QueryFilm)
		id := chi.URLParam(r, "id")
		title := chi.URLParam(r, "title")
		filter := chi.URLParam(r, "filter")

		switch {
		case id != "" && title != "": //a film
			title = strings.ReplaceAll(title, "_", " ")
			num, _ := strconv.ParseUint(id, 10, 16)
			films.Filter.ID = uint16(num)
			films.Filter.Title = title
			films.Info.SingleTitle = true
			if err := db.SearchFilm(&films.Filter, &films.Result); err != nil {
				render.Render(w, r, MyErrors.Err503)
				return
			}
		case baseURL(r.URL.Path): //all films
			films.Options.Limit = 0
			if err := db.SearchAllFilms(&films.Result); err != nil {
				render.Render(w, r, MyErrors.Err503)
				return
			}
		case filter != "": //some films
			if ok := filterFilm(films, filter); !ok {
				render.Render(w, r, MyErrors.Err400)
				return
			}
			var total int16
			if err := db.SearchSomeFilms(&films.Filter, films.Options, &films.Result, &total); err != nil {
				render.Render(w, r, MyErrors.Err503)
				return
			}
			if total > 0 {
				films.Info.TotalData = total
				films.Info.ActualPage = int16(films.Options.Offset / films.Options.Limit)
				films.Info.RemainingPages = int16((total - 1) / int16(films.Options.Limit))
				films.Info.TotalPages = films.Info.ActualPage + films.Info.RemainingPages
			}
		default:
			render.Render(w, r, MyErrors.Err404)
			return
		}

		ctx := context.WithValue(r.Context(), "films", films)
		sig.ServeHTTP(w, r.WithContext(ctx))
	})
}

// CategoryRequest Middleware that runs the inventory search to find one or
// many categories
func CategoryRequest(sig http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		categories, _ := r.Context().Value("categories").(*sakila.CategoryResult)

		switch {
		case baseURL(r.URL.Path): //all categoryes
			if err := db.SearchCategories(categories); err != nil {
				render.Render(w, r, MyErrors.Err503)
				return
			}
		default:
			render.Render(w, r, MyErrors.Err404)
			return
		}

		ctx := context.WithValue(r.Context(), "categories", categories)
		sig.ServeHTTP(w, r.WithContext(ctx))
	})
}

// ActorRequest Middleware that runs the inventory search to find one or
// many actors
func ActorRequest(sig http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		actors, _ := r.Context().Value("actors").(*sakila.ActorResult)

		switch {
		case baseURL(r.URL.Path): //all actors
			if err := db.SearchActors(actors); err != nil {
				render.Render(w, r, MyErrors.Err503)
				return
			}
		default:
			render.Render(w, r, MyErrors.Err404)
			return
		}

		ctx := context.WithValue(r.Context(), "actors", actors)
		sig.ServeHTTP(w, r.WithContext(ctx))
	})
}

func baseURL(url string) bool {
	return url == PathTo_Inventory || url == PathTo_Inventory+"/" ||
		url == PathTo_Categories || url == PathTo_Categories+"/" ||
		url == PathTo_Actors || url == PathTo_Actors+"/"
}

// Pagination Middleware that adds filters related to pagination
func Pagination(sig http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		films, _ := r.Context().Value("films").(*sakila.QueryFilm)
		if FilmPagination := chi.URLParam(r, "FilmPagination"); FilmPagination != "" {
			if ok := orderFilm(films, FilmPagination); !ok {
				render.Render(w, r, MyErrors.Err400)
				return
			}
		}

		ctx := context.WithValue(r.Context(), "films", films)
		sig.ServeHTTP(w, r.WithContext(ctx))
	})
}

func filterFilm(films *sakila.QueryFilm, rawfilter string) bool {
	filter, err := url.ParseQuery(rawfilter)
	if err != nil {
		return false
	}

	for option, val := range filter {
		qurl, err := url.PathUnescape(strings.Join(val, ""))
		if err != nil || qurl == "" {
			continue
		}

		switch option {
		case search:
			films.Filter.Description = &sakila.NullString{String: qurl, Valid: true}
		case year:
			tmp, err := strconv.ParseInt(qurl, 10, 32)
			if err != nil {
				return false
			}
			sql := sql.NullInt32{Int32: int32(tmp), Valid: true}
			films.Filter.Year = sakila.NullInt32(sql)
		case language:
			films.Filter.Lang = qurl
		case rentalDuration:
			tmp, err := strconv.ParseFloat(qurl, 32)
			if err != nil {
				return false
			}
			films.Filter.RentalDuration = float32(tmp)
		case price:
			tmp, err := strconv.ParseFloat(qurl, 32)
			if err != nil {
				return false
			}
			films.Filter.RentalPrice = float32(tmp)
		case duration:
			tmp, err := strconv.ParseInt(qurl, 10, 32)
			if err != nil {
				return false
			}
			sql := sql.NullInt32{Int32: int32(tmp), Valid: true}
			films.Filter.Lenght = sakila.NullInt32(sql)
		case replacementCost:
			tmp, err := strconv.ParseFloat(qurl, 32)
			if err != nil {
				return false
			}
			films.Filter.ReplacementCost = float32(tmp)
		case rating:
			films.Filter.Rating = &sakila.NullString{String: qurl, Valid: true}
		case specialFeatures:
			tmp := strings.ReplaceAll(qurl, "+", " ")
			arr := strings.Split(tmp, "_")
			films.Filter.SpecialFeatures = arr
		case lastUpdate:
			tmp, err := time.Parse(time.RFC3339, qurl)
			if err != nil {
				return false
			}
			sql := sql.NullTime{Time: tmp, Valid: true}
			films.Filter.LastUpdate = sakila.NullTime(sql)
		default:
			return false
		}
	}

	return len(filter) > 0
}

func orderFilm(films *sakila.QueryFilm, rawfilter string) bool {
	options, err := url.ParseQuery(rawfilter)
	if err != nil {
		return false
	}
	for option, val := range options {
		qurl, err := url.PathUnescape(strings.Join(val, ""))
		if err != nil || qurl == "" {
			continue
		}
		switch option {
		case limit:
			num, nan := strconv.ParseUint(qurl, 10, 16)
			if nan != nil {
				return false
			}
			films.Options.Limit = uint16(num)
		case offset:
			num, nan := strconv.ParseUint(qurl, 10, 16)
			if nan != nil {
				return false
			}
			films.Options.Offset = films.Options.Limit * uint16(num)
		case orderby:
			var qToArg string
			switch qurl {
			case id:
				qToArg = "id"
			case title:
				qToArg = "title"
			case year:
				qToArg = "year"
			case rentalDuration:
				qToArg = "rental duration"
			case price:
				qToArg = "price"
			case duration:
				qToArg = "duration"
			case replacementCost:
				qToArg = "replacement cost"
			case lastUpdate:
				qToArg = "last update"
			default:
				return false
			}
			films.Options.OrderBy = qToArg
		case order:
			films.Options.Order = qurl
		default:
			return false
		}
	}
	return true
}
