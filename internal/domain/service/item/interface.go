package item

import "github.com/omcrgnt/demo/internal/domain/model"

type ItemRepository interface {
	AllItems() ([]model.Item, error)
	ItemByID(id string) (model.Item, error)
	AddItem(title string) (model.Item, error)
	SetItem(id, title string) (model.Item, error)
	RemoveItem(id string) error
}
