package testdata

import (
	"context"
	"errors"
)

//go:generate go run ../cmd/obsgen/main.go -type=Repo
type Repo struct {
	name string
}

func NewRepo(name string) *Repo {
	return &Repo{
		name: name,
	}
}

func (r *Repo) Label() string {
	return r.name
}

func (r *Repo) Save(ctx context.Context, data string) error {
	if data == "" {
		return errors.New("empty data")
	}
	return nil
}
