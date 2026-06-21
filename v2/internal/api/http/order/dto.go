package order

import "github.com/omcrgnt/demo/v2/internal/domain/model"

type OrderResponse struct {
	ID    string `json:"id"`
	Title string `json:"title"`
}

type CreateOrderRequest struct {
	Title string `json:"title"`
}

type UpdateOrderRequest struct {
	Title string `json:"title"`
}

func toOrderResponse(order model.Order) OrderResponse {
	return OrderResponse{
		ID:    order.ID,
		Title: order.Title,
	}
}

func toOrderResponses(orders []model.Order) []OrderResponse {
	out := make([]OrderResponse, len(orders))
	for i, order := range orders {
		out[i] = toOrderResponse(order)
	}
	return out
}
