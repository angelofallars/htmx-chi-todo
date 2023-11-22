package main

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/angelofallars/htmx-chi-todo/todo"
	"github.com/angelofallars/htmx-chi-todo/user"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/jwtauth/v5"
	_ "github.com/mattn/go-sqlite3"
	"github.com/redis/go-redis/v9"
)

const sqliteFileName = "sqlite.db"

func main() {
	redisClient := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
		DB:   0,
	})

	sqliteDB, err := sql.Open("sqlite3", sqliteFileName)
	if err != nil {
		log.Fatal(err)
	}

	userSQLite3Repo := user.NewSQLiteRepository(sqliteDB)

	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	user.NewHandler(
		user.NewService(
			userSQLite3Repo,
		),
	).Mount(r)

	tokenAuth := jwtauth.New("HS256", []byte("secret"), nil)

	r.Group(func(r chi.Router) {
		r.Use(jwtauth.Verifier(tokenAuth))
		r.Use(jwtauth.Authenticator)

		todo.NewHandler(
			todo.NewService(
				todo.NewRedisRepository(redisClient),
			),
		).Mount(r)
	})

	r.Handle("/assets/*", http.StripPrefix("/assets/", http.FileServer(http.Dir("assets"))))

	log.Fatal(http.ListenAndServe(":3000", r))
}
