package item

import "github.com/omcrgnt/demo/internal/domain/model"

type ItemResponse struct {
	ID    string `json:"id"`
	Title string `json:"title"`
}

type CreateItemRequest struct {
	Title string `json:"title"`
}

type UpdateItemRequest struct {
	Title string `json:"title"`
}

func toItemResponse(item model.Item) ItemResponse {
	return ItemResponse{
		ID:    item.ID,
		Title: item.Title,
	}
}

func toItemResponses(items []model.Item) []ItemResponse {
	out := make([]ItemResponse, len(items))
	for i, item := range items {
		out[i] = toItemResponse(item)
	}
	return out
}
