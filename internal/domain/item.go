package domain

type Item struct {
	ID    string `json:"id"`
	Title string `json:"title"`
}

type CreateItemRequest struct {
	Title string `json:"title"`
}

type UpdateItemRequest struct {
	Title string `json:"title"`
}
