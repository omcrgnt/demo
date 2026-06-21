package runner

import (
	"context"

	"github.com/omcrgnt/demo/v2/pkg/res"
)

// Pool is a resource source for Run/Stop; [res.Registry] satisfies it.
type Pool interface {
	WalkEntries(fn func(res.Entry) bool)
}

type Starter interface {
	Start(ctx context.Context) error
}

type Closer interface {
	Close(ctx context.Context) error
}

type StartCloser interface {
	Starter
	Closer
}
