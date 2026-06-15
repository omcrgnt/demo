package http

import (
	"encoding/json"
	"errors"
	nethttp "net/http"

	"github.com/go-chi/chi/v5"
	"github.com/omcrgnt/demo/internal/domain"
	"github.com/omcrgnt/demo/internal/domain/service/item"
)

type API struct {
	mux *chi.Mux
	svc *item.Service
}

type Config struct{}

func (Config) Build() (any, error) {
	return &API{mux: chi.NewRouter()}, nil
}

func (a *API) Deps() []any {
	return []any{
		(*item.Service)(nil),
	}
}

func (a *API) Inject(args []any) {
	for _, arg := range args {
		if svc, ok := arg.(*item.Service); ok {
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
		writeError(w, nethttp.StatusInternalServerError, err)
		return
	}
	if items == nil {
		writeJSON(w, nethttp.StatusOK, []ItemResponse{})
		return
	}
	writeJSON(w, nethttp.StatusOK, toItemResponses(items))
}

func (a *API) getItem(w nethttp.ResponseWriter, r *nethttp.Request) {
	id := chi.URLParam(r, "id")
	item, err := a.svc.Get(r.Context(), id)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	writeJSON(w, nethttp.StatusOK, toItemResponse(item))
}

func (a *API) createItem(w nethttp.ResponseWriter, r *nethttp.Request) {
	var req CreateItemRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, nethttp.StatusBadRequest, err)
		return
	}

	item, err := a.svc.Create(r.Context(), req.Title)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	writeJSON(w, nethttp.StatusCreated, toItemResponse(item))
}

func (a *API) updateItem(w nethttp.ResponseWriter, r *nethttp.Request) {
	id := chi.URLParam(r, "id")
	var req UpdateItemRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, nethttp.StatusBadRequest, err)
		return
	}

	item, err := a.svc.Update(r.Context(), id, req.Title)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	writeJSON(w, nethttp.StatusOK, toItemResponse(item))
}

func (a *API) deleteItem(w nethttp.ResponseWriter, r *nethttp.Request) {
	id := chi.URLParam(r, "id")
	if err := a.svc.Delete(r.Context(), id); err != nil {
		writeServiceError(w, err)
		return
	}
	w.WriteHeader(nethttp.StatusNoContent)
}

func writeServiceError(w nethttp.ResponseWriter, err error) {
	switch {
	case errors.Is(err, domain.ErrNotFound):
		writeError(w, nethttp.StatusNotFound, err)
	case errors.Is(err, domain.ErrInvalidInput):
		writeError(w, nethttp.StatusBadRequest, err)
	default:
		writeError(w, nethttp.StatusInternalServerError, err)
	}
}

func writeError(w nethttp.ResponseWriter, status int, err error) {
	writeJSON(w, status, map[string]string{"error": err.Error()})
}

func writeJSON(w nethttp.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}
