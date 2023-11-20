package main

import (
	"net/http"

	"github.com/angelofallars/htmx-chi-todo/todo"
	"github.com/angelofallars/htmx-chi-todo/user"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/jwtauth/v5"
	"github.com/redis/go-redis/v9"
	"github.com/unrolled/render"
)

func main() {
	redisClient := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
		DB:   0,
	})

	rnd := render.New()

	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	user.NewHandler(
		user.NewService(
			user.NewRedisRepository(redisClient),
		),
	).Mount(r)

	tokenAuth := jwtauth.New("HS256", []byte("secret"), nil)

	r.Group(func(r chi.Router) {
		r.Use(jwtauth.Verifier(tokenAuth))
		r.Use(jwtauth.Authenticator)

		todo.NewHandler(
			rnd,
			todo.NewService(
				todo.NewRedisRepository(redisClient),
			),
		).Mount(r)
	})

	r.Handle("/assets/*", http.StripPrefix("/assets/", http.FileServer(http.Dir("assets"))))

	http.ListenAndServe(":3000", r)
}
