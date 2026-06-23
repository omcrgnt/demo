package http

import (
	"encoding/json"
	nethttp "net/http"

	"github.com/go-chi/chi/v5"
	"github.com/omcrgnt/demo/internal/domain/service/item"
)

type API struct {
	mux *chi.Mux
	svc item.ItemService
}

func (a *API) NewResource() (any, error) {
	return &API{mux: chi.NewRouter()}, nil
}

func (a *API) Deps() []any {
	return []any{
		(*item.ItemService)(nil),
	}
}

func (a *API) Inject(args []any) {
	for _, arg := range args {
		if svc, ok := arg.(item.ItemService); ok {
			a.svc = svc
		}
	}
	a.registerRoutes()
}

func (a *API) ServeHTTP(w nethttp.ResponseWriter, r *nethttp.Request) {
	a.mux.ServeHTTP(w, r)
}

func (a *API) registerRoutes() {
	a.mux.Get("/items", a.listItems)
	a.mux.Get("/items/{id}", a.getItem)
	a.mux.Post("/items", a.createItem)
	a.mux.Put("/items/{id}", a.updateItem)
	a.mux.Delete("/items/{id}", a.deleteItem)
}

func (a *API) listItems(w nethttp.ResponseWriter, r *nethttp.Request) {
	items, err := a.svc.List(r.Context())
	if err != nil {
		WriteError(w, nethttp.StatusInternalServerError, err)
		return
	}
	if items == nil {
		WriteJSON(w, nethttp.StatusOK, []ItemResponse{})
		return
	}
	WriteJSON(w, nethttp.StatusOK, toItemResponses(items))
}

func (a *API) getItem(w nethttp.ResponseWriter, r *nethttp.Request) {
	id := chi.URLParam(r, "id")
	item, err := a.svc.Get(r.Context(), id)
	if err != nil {
		WriteServiceError(w, err)
		return
	}
	WriteJSON(w, nethttp.StatusOK, toItemResponse(item))
}

func (a *API) createItem(w nethttp.ResponseWriter, r *nethttp.Request) {
	var req CreateItemRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteError(w, nethttp.StatusBadRequest, err)
		return
	}

	item, err := a.svc.Create(r.Context(), req.Title)
	if err != nil {
		WriteServiceError(w, err)
		return
	}
	WriteJSON(w, nethttp.StatusCreated, toItemResponse(item))
}

func (a *API) updateItem(w nethttp.ResponseWriter, r *nethttp.Request) {
	id := chi.URLParam(r, "id")
	var req UpdateItemRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteError(w, nethttp.StatusBadRequest, err)
		return
	}

	item, err := a.svc.Update(r.Context(), id, req.Title)
	if err != nil {
		WriteServiceError(w, err)
		return
	}
	WriteJSON(w, nethttp.StatusOK, toItemResponse(item))
}

func (a *API) deleteItem(w nethttp.ResponseWriter, r *nethttp.Request) {
	id := chi.URLParam(r, "id")
	if err := a.svc.Delete(r.Context(), id); err != nil {
		WriteServiceError(w, err)
		return
	}
	w.WriteHeader(nethttp.StatusNoContent)
}
