package order

import (
	"encoding/json"
	nethttp "net/http"

	"github.com/go-chi/chi/v5"
	"github.com/omcrgnt/demo/v2/internal/api/http"
	ordersvc "github.com/omcrgnt/demo/v2/internal/domain/service/order"
)

type API struct {
	mux *chi.Mux
	svc ordersvc.OrderService
}

func (a *API) NewResource() (any, error) {
	return &API{mux: chi.NewRouter()}, nil
}

func (a *API) Deps() []any {
	return []any{
		(*ordersvc.OrderService)(nil),
	}
}

func (a *API) Inject(args []any) {
	for _, arg := range args {
		if svc, ok := arg.(ordersvc.OrderService); ok {
			a.svc = svc
		}
	}
	a.registerRoutes()
}

func (a *API) ServeHTTP(w nethttp.ResponseWriter, r *nethttp.Request) {
	a.mux.ServeHTTP(w, r)
}

func (a *API) registerRoutes() {
	a.mux.Get("/orders", a.listOrders)
	a.mux.Get("/orders/{id}", a.getOrder)
	a.mux.Post("/orders", a.createOrder)
	a.mux.Put("/orders/{id}", a.updateOrder)
	a.mux.Delete("/orders/{id}", a.deleteOrder)
}

func (a *API) listOrders(w nethttp.ResponseWriter, r *nethttp.Request) {
	orders, err := a.svc.List(r.Context())
	if err != nil {
		http.WriteError(w, nethttp.StatusInternalServerError, err)
		return
	}
	if orders == nil {
		http.WriteJSON(w, nethttp.StatusOK, []OrderResponse{})
		return
	}
	http.WriteJSON(w, nethttp.StatusOK, toOrderResponses(orders))
}

func (a *API) getOrder(w nethttp.ResponseWriter, r *nethttp.Request) {
	id := chi.URLParam(r, "id")
	order, err := a.svc.Get(r.Context(), id)
	if err != nil {
		http.WriteServiceError(w, err)
		return
	}
	http.WriteJSON(w, nethttp.StatusOK, toOrderResponse(order))
}

func (a *API) createOrder(w nethttp.ResponseWriter, r *nethttp.Request) {
	var req CreateOrderRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.WriteError(w, nethttp.StatusBadRequest, err)
		return
	}

	order, err := a.svc.Create(r.Context(), req.Title)
	if err != nil {
		http.WriteServiceError(w, err)
		return
	}
	http.WriteJSON(w, nethttp.StatusCreated, toOrderResponse(order))
}

func (a *API) updateOrder(w nethttp.ResponseWriter, r *nethttp.Request) {
	id := chi.URLParam(r, "id")
	var req UpdateOrderRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.WriteError(w, nethttp.StatusBadRequest, err)
		return
	}

	order, err := a.svc.Update(r.Context(), id, req.Title)
	if err != nil {
		http.WriteServiceError(w, err)
		return
	}
	http.WriteJSON(w, nethttp.StatusOK, toOrderResponse(order))
}

func (a *API) deleteOrder(w nethttp.ResponseWriter, r *nethttp.Request) {
	id := chi.URLParam(r, "id")
	if err := a.svc.Delete(r.Context(), id); err != nil {
		http.WriteServiceError(w, err)
		return
	}
	w.WriteHeader(nethttp.StatusNoContent)
}
