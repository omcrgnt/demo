package test

import (
	"context"
	"errors"
)

//go:generate obsgen -type=Repo
type Repo struct {
	Name string
}

func (r *Repo) Label() string {
	return r.Name
}

func (r *Repo) Save(ctx context.Context, data string) error {
	if data == "" {
		return errors.New("database error")
	}
	return nil
}
