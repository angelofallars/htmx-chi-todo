package user

import (
	"net/http"

	"github.com/angelofallars/htmx-chi-todo/htmx"
	svc "github.com/angelofallars/htmx-chi-todo/service"
	"github.com/angelofallars/htmx-chi-todo/site"
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

func (h handler) Signup(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	un := r.Form.Get("username")
	ea := r.Form.Get("email")
	pw := r.Form.Get("password")

	username, err := newUsername(un)
	if err != nil {
		site.RenderError(w,
			http.StatusBadRequest,
			err,
		)
		return
	}

	password, err := newHashedPassword(pw)
	if err != nil {
		site.RenderError(w,
			http.StatusBadRequest,
			err,
		)
		return
	}

	email, err := newEmail(ea)
	if err != nil {
		site.RenderError(w,
			http.StatusBadRequest,
			err,
		)
		return
	}

	user := newUser(username, password, email)

	err = h.service.Signup(r.Context(), user)
	if err != nil {
		site.RenderError(w,
			http.StatusInternalServerError,
			err,
		)
		return
	}

	w.Header().Add(htmx.HeaderRedirect, "/login")
}

func (h handler) LoginPage(w http.ResponseWriter, r *http.Request) {
	site.RenderRootOrPartial(w, r,
		"Login",
		loginPage(),
	)
}

func (h handler) Login(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	un := r.Form.Get("username")
	pw := r.Form.Get("password")

	username, err := newUsername(un)
	if err != nil {
		site.RenderError(w,
			http.StatusBadRequest,
			err,
		)
		return
	}

	jwt, err := h.service.Login(r.Context(), username, pw)
	if err != nil {
		site.RenderError(w,
			http.StatusInternalServerError,
			err,
		)
		return
	}

	w.Header().Add(htmx.HeaderRedirect, "/login")
	http.SetCookie(w, &http.Cookie{
		Name:  "jwt",
		Value: jwt,
	})
	w.Write([]byte(jwt))
}
