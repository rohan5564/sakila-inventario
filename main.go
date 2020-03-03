package main

import (
	"github.com/rohan5564/sakila-inventario/sakila"
	"github.com/rohan5564/sakila-inventario/sakila/rest"

	"context"
	"database/sql"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
	_ "github.com/go-sql-driver/mysql"
)

var port uint

var session sakila.DBdata

func init() {
	def := sakila.DBdata{
		User:     "jean",
		Pass:     "64fa4632",
		Schema:   "sakila",
		Protocol: "tcp",
		IP:       "localhost",
		Port:     3306}

	flag.StringVar(&session.User, "user", def.User, "Database user")
	flag.StringVar(&session.Pass, "pass", def.Pass, "Database password")
	flag.StringVar(&session.Schema, "schema", def.Schema, "Database schema")
	flag.StringVar(&session.Protocol, "protocol", def.Protocol, "Database protocol connection")
	flag.StringVar(&session.IP, "ip", def.IP, "Database IP")
	flag.UintVar(&session.Port, "port", def.Port, "Database port")
	flag.UintVar(&port, "server-port", 80, "port where API will run")

}

func main() {
	flag.Parse()

	exitCode := func(code int) int {
		sec := time.Tick(time.Second)
		for i := 3; i >= 0; i-- {
			<-sec
			fmt.Printf("\rClosing in %d", i)
		}
		return code
	}

	dburl := fmt.Sprintf("%s:%s@%s(%s:%d)/%s?timeout=5s&parseTime=true&multiStatements=true",
		session.User, session.Pass, session.Protocol, session.IP, session.Port, session.Schema)

	err := StartDB(dburl)
	defer sakila.Statement.Close()

	if err != nil {
		os.Exit(exitCode(1))
	}

	fmt.Println("Generating routes...")
	router := API(middleware.Logger, middleware.Recoverer, middleware.StripSlashes)
	server := &http.Server{Addr: fmt.Sprintf(":%d", port), Handler: router}

	go func() {
		log.Println("Server running on port ", port)
		if err := server.ListenAndServe(); err != http.ErrServerClosed {
			log.Println(err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)

	<-stop

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Println(err)
	}

	os.Exit(exitCode(0))
}

// StartDB Opens the database connection
func StartDB(dburl string) error {
	var err error
	sakila.Statement, err = sql.Open("mysql", dburl)
	if err != nil {
		log.Println(err)
		return err
	}

	if err := sakila.Statement.Ping(); err != nil {
		fmt.Println("connection refused: " + dburl)
		return err
	}

	return nil
}

// API Loads the api path and returns it as a new router
func API(middlewares ...func(http.Handler) http.Handler) *chi.Mux {
	router := chi.NewRouter()
	for _, mid := range middlewares {
		router.Use(mid)
	}
	router.Use(render.SetContentType(render.ContentTypeJSON))

	router.Route(rest.PathTo_Inventory, func(router chi.Router) {
		router.Use(rest.NewFilm)
		router.With(rest.FilmRequest).Get("/", rest.GetFilm)
		router.With(rest.FilmRequest).Post("/", rest.CreateFilm)
		router.Route("/{FilmPagination}", func(router chi.Router) {
			router.Use(rest.Pagination, rest.FilmRequest)
			router.Get("/", rest.GetFilm)
		})

		router.Route("/search/{filter}", func(router chi.Router) {
			router.With(rest.FilmRequest).Get("/", rest.GetFilm)
			router.Group(func(router chi.Router) {
				router.Use(rest.Pagination, rest.FilmRequest)
				router.Get("/{FilmPagination}", rest.GetFilm)
			})
		})

		// github.com/google/re2/wiki/Syntax
		router.Route("/film/{id:\\d+}/{title:[[:alnum:]\\p{Latin}%_]+}", func(router chi.Router) {
			router.Use(rest.FilmRequest)
			router.Get("/", rest.GetFilm)
			router.Put("/", rest.UpdateFilm)
			router.Delete("/", rest.DeleteFilm)
		})
	})

	router.Route(rest.PathTo_Categories, func(router chi.Router) {
		router.Use(rest.NewCategory)
		router.Use(rest.CategoryRequest)
		router.Get("/", rest.GetCategories)
	})

	router.Route(rest.PathTo_Actors, func(router chi.Router) {
		router.Use(rest.NewActor)
		router.Use(rest.ActorRequest)
		router.Get("/", rest.GetActors)
	})

	return router
}
