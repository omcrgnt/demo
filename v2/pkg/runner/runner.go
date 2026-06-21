package runner

import (
	"context"
	"errors"
	"fmt"

	"github.com/omcrgnt/demo/v2/pkg/res"
	"golang.org/x/sync/errgroup"
)

// Runner runs and stops [Starter]/[Closer] resources from a [Pool].
type Runner struct{}

func (r *Runner) NewResource() (any, error) {
	return &Runner{}, nil
}

type engine struct {
	resources []any
}

func newEngine(pool Pool) *engine {
	return &engine{resources: collect(pool)}
}

// New builds a runner from pool without registering a resource (tests).
func New(pool Pool) *engine {
	return newEngine(pool)
}

func collect(pool Pool) []any {
	var resources []any
	pool.WalkEntries(func(e res.Entry) bool {
		resources = append(resources, e.Value)
		return true
	})
	return resources
}

func (r *Runner) Run(ctx context.Context, pool Pool) error {
	return newEngine(pool).run(ctx)
}

func (r *Runner) Stop(ctx context.Context, pool Pool) error {
	return newEngine(pool).stop(ctx)
}

func (e *engine) run(rctx context.Context) error {
	group, ctx := errgroup.WithContext(rctx)

	for _, res := range e.resources {
		if s, ok := res.(Starter); ok {
			starter := s
			group.Go(func() error {
				if err := starter.Start(ctx); err != nil {
					return fmt.Errorf("starter %T failed: %w", starter, err)
				}
				return nil
			})
		}
	}

	return group.Wait()
}

func (e *engine) stop(ctx context.Context) error {
	var errs []error

	for i := len(e.resources) - 1; i >= 0; i-- {
		res := e.resources[i]
		if sc, ok := res.(StartCloser); ok {
			if err := sc.Close(ctx); err != nil {
				errs = append(errs, fmt.Errorf("active resource %T: %w", res, err))
			}
		}
	}

	for i := len(e.resources) - 1; i >= 0; i-- {
		res := e.resources[i]

		closer, isCloser := res.(Closer)
		_, isActive := res.(StartCloser)

		if isCloser && !isActive {
			if err := closer.Close(ctx); err != nil {
				errs = append(errs, fmt.Errorf("passive resource %T: %w", res, err))
			}
		}
	}
	return errors.Join(errs...)
}
