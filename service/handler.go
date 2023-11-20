package service

import "github.com/go-chi/chi/v5"

type HandlerMounter interface {
	// Mount this handler's endpoints into a router.
	Mount(r chi.Router)
}
