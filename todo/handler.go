package todo

import (
	"errors"
	"fmt"
	"net/http"

	svc "github.com/angelofallars/htmx-chi-todo/service"
	"github.com/angelofallars/htmx-chi-todo/site"
	"github.com/angelofallars/htmx-go"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type (
	Handler interface {
		svc.HandlerMounter
		Page(w http.ResponseWriter, r *http.Request)
		GetList(w http.ResponseWriter, r *http.Request)
		GetItem(w http.ResponseWriter, r *http.Request)
		CreateItem(w http.ResponseWriter, r *http.Request)
		ToggleItemComplete(w http.ResponseWriter, r *http.Request)
		DeleteItem(w http.ResponseWriter, r *http.Request)
	}
	handler struct {
		service Service
	}
)

func (h handler) Mount(r chi.Router) {
	r.Get("/", h.Page)
	r.Get("/items", h.GetList)
	r.Get("/items/{id}", h.GetItem)
	r.Post("/items", h.CreateItem)
	r.Put("/items/{id}/toggle", h.ToggleItemComplete)
	r.Put("/items/{id}", h.UpdateItem)
	r.Delete("/items/{id}", h.DeleteItem)
}

func NewHandler(service Service) Handler {
	return &handler{
		service: service,
	}
}

func (h handler) Page(w http.ResponseWriter, r *http.Request) {
	list := newList()

	item1 := newItem("Write notes on Chemistry 1", "")
	item2 := newItem("Listen to CS lecture 240", "Take notes on Data structures & Algorithms")
	item2.IsDone = true
	item3 := newItem("Study HTMX", "HTMX is the best!")
	item4 := newItem("Finish Chapter 12 of the Rust book", "")

	list.Items = []*item{
		item1,
		item2,
		item3,
		item4,
	}

	_, err := h.service.CreateList(r.Context(), list)
	_, err = h.service.CreateItem(r.Context(), item1)
	_, err = h.service.CreateItem(r.Context(), item2)
	_, err = h.service.CreateItem(r.Context(), item3)
	_, err = h.service.CreateItem(r.Context(), item4)
	if err != nil {
		site.RenderError(w,
			http.StatusInternalServerError,
			err,
		)
		return
	}

	w.Header().Add(htmx.HeaderReplaceUrl, fmt.Sprintf("/%v", list.ID.String()))

	site.RenderRootOrPartial(w, r,
		"TODO",
		page(list),
	)
}

func (h handler) CreateItem(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	name := r.Form.Get("task-name")
	description := r.Form.Get("task-description")

	if len(name) == 0 {
		site.RenderError(w,
			http.StatusBadRequest,
			errors.New("Cannot have task names with a length of zero"),
		)
		return
	}

	item := newItem(name, description)

	_, err := h.service.CreateItem(r.Context(), item)

	if err != nil {
		site.RenderError(w,
			http.StatusInternalServerError,
			err,
		)
		return
	}

	item.Component().Render(r.Context(), w)
}

func (h handler) GetList(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		site.RenderError(w, http.StatusBadRequest, err)
		return
	}

	list, err := h.service.GetList(r.Context(), &id)
	if err != nil {
		site.RenderError(w,
			http.StatusInternalServerError,
			err,
		)
		return
	}

	list.Component().Render(r.Context(), w)
}

func (h handler) GetItem(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		site.RenderError(w, http.StatusBadRequest, err)
		return
	}

	item, err := h.service.GetItem(r.Context(), &id)
	if err != nil {
		site.RenderError(w, http.StatusInternalServerError, err)
		return
	}

	item.Component().Render(r.Context(), w)
}

func (h handler) UpdateItem(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		site.RenderError(w, http.StatusInternalServerError, err)
		return
	}

	r.ParseForm()
	name := r.Form.Get("task-name")
	description := r.Form.Get("task-description")

	item, err := h.service.UpdateItem(r.Context(), &id, name, description)
	if err != nil {
		site.RenderError(w,
			http.StatusInternalServerError,
			err,
		)
		return
	}

	item.Component().Render(r.Context(), w)
}

func (h handler) ToggleItemComplete(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		site.RenderError(w, http.StatusBadRequest, err)
		return
	}

	completionStatus, err := h.service.ToggleItemComplete(r.Context(), &id)
	if err != nil {
		site.RenderError(w, http.StatusInternalServerError, err)
		return
	}

	completionStatus.Component(id).Render(r.Context(), w)
}

func (h handler) DeleteItem(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		site.RenderError(w, http.StatusInternalServerError, err)
		return
	}

	err = h.service.DeleteItem(r.Context(), &id)
	if err != nil {
		site.RenderError(w, http.StatusInternalServerError, err)
		return
	}
}
