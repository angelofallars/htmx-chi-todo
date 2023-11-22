package user

import (
	"errors"
	"net/http"

	"github.com/angelofallars/htmx-chi-todo/service"
	svc "github.com/angelofallars/htmx-chi-todo/service"
	"github.com/angelofallars/htmx-chi-todo/site"
	"github.com/angelofallars/htmx-go"
	"github.com/go-chi/chi/v5"
)

type (
	Handler interface {
		svc.HandlerMounter
		SignupPage(w http.ResponseWriter, r *http.Request)
		Signup(w http.ResponseWriter, r *http.Request)
		LoginPage(w http.ResponseWriter, r *http.Request)
		Login(w http.ResponseWriter, r *http.Request)
	}

	handler struct {
		service Service
	}
)

func (h handler) Mount(r chi.Router) {
	r.Get("/signup", h.SignupPage)
	r.Post("/signup", h.Signup)
	r.Get("/login", h.LoginPage)
	r.Post("/login", h.Login)
}

func NewHandler(service Service) Handler {
	return &handler{
		service: service,
	}
}

func (h handler) SignupPage(w http.ResponseWriter, r *http.Request) {
	site.RenderRootOrPartial(w, r,
		"Sign Up",
		signupPage(),
	)
}

func (h handler) LoginPage(w http.ResponseWriter, r *http.Request) {
	site.RenderRootOrPartial(w, r,
		"Login",
		loginPage(),
	)
}

func (h handler) Signup(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	username := r.Form.Get("username")
	email := r.Form.Get("email")
	password := r.Form.Get("password")

	err := h.service.Signup(r.Context(), SignupReq{
		Username:    username,
		Email:       email,
		RawPassword: password,
	})

	if err != nil {
		code := http.StatusInternalServerError
		if errors.Is(err, service.ErrValidation) {
			code = http.StatusBadRequest
		}
		site.RenderError(w,
			code,
			err,
		)
		return
	}

	htmx.NewResponse().
		Redirect("login").
		Write(w)
}

func (h handler) Login(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	username := r.Form.Get("username")
	password := r.Form.Get("password")

	jwt, err := h.service.Login(r.Context(), LoginReq{
		Username:    username,
		RawPassword: password,
	})

	if err != nil {
		code := http.StatusInternalServerError
		if errors.Is(err, service.ErrValidation) {
			code = http.StatusBadRequest
		}
		site.RenderError(w,
			code,
			err,
		)
		return
	}

	htmx.NewResponse().
		Redirect("/").
		Write(w)

	http.SetCookie(w, &http.Cookie{
		Name:  "jwt",
		Value: jwt,
	})
}
