package httpapi

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/omcrgnt/demo/internal/domain"
	"github.com/omcrgnt/demo/internal/service"
)

type ItemService interface {
	List() ([]domain.Item, error)
	Get(id string) (domain.Item, error)
	Create(title string) (domain.Item, error)
	Update(id, title string) (domain.Item, error)
	Delete(id string) error
}

type API struct {
	mux *chi.Mux
	svc ItemService
}

type Config struct{}

func (Config) Build() (any, error) {
	return &API{mux: chi.NewRouter()}, nil
}

func (a *API) Deps() []any {
	return []any{
		(*ItemService)(nil),
	}
}

func (a *API) Inject(args []any) {
	for _, arg := range args {
		if svc, ok := arg.(ItemService); ok {
			a.svc = svc
		}
	}
	a.registerRoutes()
}

func (a *API) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	a.mux.ServeHTTP(w, r)
}

func (a *API) registerRoutes() {
	a.mux.Get("/items", a.listItems)
	a.mux.Get("/items/{id}", a.getItem)
	a.mux.Post("/items", a.createItem)
	a.mux.Put("/items/{id}", a.updateItem)
	a.mux.Delete("/items/{id}", a.deleteItem)
}

func (a *API) listItems(w http.ResponseWriter, r *http.Request) {
	items, err := a.svc.List()
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}
	if items == nil {
		items = []domain.Item{}
	}
	writeJSON(w, http.StatusOK, items)
}

func (a *API) getItem(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	item, err := a.svc.Get(id)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, item)
}

func (a *API) createItem(w http.ResponseWriter, r *http.Request) {
	var req domain.CreateItemRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	item, err := a.svc.Create(req.Title)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, item)
}

func (a *API) updateItem(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var req domain.UpdateItemRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	item, err := a.svc.Update(id, req.Title)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, item)
}

func (a *API) deleteItem(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if err := a.svc.Delete(id); err != nil {
		writeServiceError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func writeServiceError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, service.ErrNotFound):
		writeError(w, http.StatusNotFound, err)
	case errors.Is(err, service.ErrInvalidInput):
		writeError(w, http.StatusBadRequest, err)
	default:
		writeError(w, http.StatusInternalServerError, err)
	}
}

func writeError(w http.ResponseWriter, status int, err error) {
	writeJSON(w, status, map[string]string{"error": err.Error()})
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}
