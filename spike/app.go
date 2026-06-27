package main

import (
	"fmt"
	"os"
	"reflect"
	"strings"

	"github.com/omcrgnt/ecfg"
	_ "github.com/omcrgnt/logger/use"
	_ "github.com/omcrgnt/telemetry/use"

	"github.com/omcrgnt/demo/internal/api/http"
	orderhttp "github.com/omcrgnt/demo/internal/api/http/order"
	"github.com/omcrgnt/demo/internal/data/sync/item/memory"
	ordermemory "github.com/omcrgnt/demo/internal/data/sync/order/memory"
	"github.com/omcrgnt/demo/internal/domain/service/item"
	"github.com/omcrgnt/demo/internal/domain/service/order"
	"github.com/omcrgnt/res"
	"github.com/omcrgnt/res/unique"
	"github.com/omcrgnt/runner"
	srvhttp "github.com/omcrgnt/srv-http"
)

type DecoratorFunc func(any) any

func validateMainT01(reg *unique.Registry) error {
	return walkPool(reg, func(v any) error {
		if _, ok := v.(Spec); ok {
			return fmt.Errorf("T01: spec %T on main", v)
		}
		if isResourcerDeferred(v) {
			return fmt.Errorf("T01: deferred resourcer %T on main", v)
		}
		return nil
	})
}

func validateSideT11(reg *unique.Registry) error {
	return walkPool(reg, func(v any) error {
		if isSpec(v) {
			return nil
		}
		if isResourcerDeferred(v) {
			return fmt.Errorf("T11: unexpected deferred resourcer %T", v)
		}
		return nil
	})
}

func walkPool(reg *unique.Registry, fn func(any) error) error {
	var err error
	reg.WalkEntries(func(e res.Entry) bool {
		if err = fn(e.Value); err != nil {
			return false
		}
		return true
	})
	return err
}

func isSpec(v any) bool {
	_, ok := v.(Spec)
	return ok
}

type resourcerDeferred struct {
	factory Resourcer
}

func isResourcerDeferred(v any) bool {
	_, ok := v.(resourcerDeferred)
	return ok
}

func t11FillSide(side *unique.Registry, appResources any, prefix string) error {
	rv, err := structValue(appResources)
	if err != nil {
		return err
	}
	rt := rv.Type()

	for i := 0; i < rv.NumField(); i++ {
		sf := rt.Field(i)
		fv := rv.Field(i)
		if !fv.CanInterface() {
			continue
		}
		field := fv.Interface()

		cfg, isCfg := field.(Configurable)
		resr, isRes := field.(Resourcer)
		switch {
		case isCfg && isRes:
			return fmt.Errorf("T11: %s: both Configurable and Resourcer", sf.Name)
		case isCfg:
			spec, err := cfg.BuildConfig()
			if err != nil {
				return fmt.Errorf("T11: %s: BuildConfig: %w", sf.Name, err)
			}
			if err := applyEnvSpike(spec, sf, prefix); err != nil {
				return err
			}
			if err := side.Add(spec); err != nil {
				return fmt.Errorf("T11: %s: %w", sf.Name, err)
			}
		case isRes:
			built, err := resr.NewResource()
			if err != nil {
				return fmt.Errorf("T11: %s: NewResource: %w", sf.Name, err)
			}
			if err := side.Add(built); err != nil {
				return fmt.Errorf("T11: %s: %w", sf.Name, err)
			}
		default:
			return fmt.Errorf("T11: %s: must be Configurable or Resourcer", sf.Name)
		}
	}
	return validateSideT11(side)
}

func applyEnvSpike(spec Spec, sf reflect.StructField, prefix string) error {
	_ = spec
	_ = sf
	_ = prefix
	return nil
}

func ecfgKey(prefix, block, field string) string {
	p := strings.TrimSuffix(prefix, "_")
	if p != "" {
		p += "_"
	}
	return p + block + "_" + field
}

func t12Materialize(side *unique.Registry) error {
	var specs []Spec
	side.WalkEntries(func(e res.Entry) bool {
		if s, ok := e.Value.(Spec); ok {
			specs = append(specs, s)
		}
		return true
	})
	for _, x := range specs {
		built, err := x.Build()
		if err != nil {
			return fmt.Errorf("T12: %T: %w", x, err)
		}
		if err := side.Remove(x); err != nil {
			return err
		}
		if err := side.Add(built); err != nil {
			return fmt.Errorf("T12: add %T: %w", built, err)
		}
	}
	return nil
}

func t30Transform(reg *unique.Registry, fns ...DecoratorFunc) error {
	if len(fns) == 0 {
		return nil
	}
	wrapped := make([]res.TransformFunc, len(fns))
	for i, fn := range fns {
		wrapped[i] = res.TransformFunc(fn)
	}
	return reg.Transform(wrapped...)
}

func bootstrap(appResources any, transforms ...DecoratorFunc) (*unique.Registry, error) {
	registryMain := unique.Global()

	if err := validateMainT01(registryMain); err != nil {
		return nil, err
	}

	registrySide := unique.New()

	if err := t11FillSide(registrySide, appResources, envPrefix); err != nil {
		return nil, err
	}

	if err := t12Materialize(registrySide); err != nil {
		return nil, err
	}

	if err := unique.Merge(registryMain, registrySide); err != nil {
		return nil, err
	}

	if err := t30Transform(registryMain, transforms...); err != nil {
		return nil, err
	}

	return registryMain, nil
}

func structValue(v any) (reflect.Value, error) {
	if v == nil {
		return reflect.Value{}, fmt.Errorf("nil app resources")
	}
	rv := reflect.ValueOf(v)
	for rv.Kind() == reflect.Ptr {
		if rv.IsNil() {
			return reflect.Value{}, fmt.Errorf("nil app resources")
		}
		rv = rv.Elem()
	}
	if rv.Kind() != reflect.Struct {
		return reflect.Value{}, fmt.Errorf("want struct, got %s", rv.Kind())
	}
	return rv, nil
}

// main
const envPrefix = "DEMO"

var appResources AppResources

type AppResources struct {
	Runner          *runner.Runner
	RepoItem        memory.RepoRoot
	ServiceItem     *item.Service `ecfg:"SERVICE_ITEM"`
	RepoOrder       ordermemory.RepoRoot
	ServiceOrder    *order.Service
	Metrics         http.Metrics
	ServerHTTPItem  *http.Server `ecfg:"SERVER_HTTP_ITEM"`
	APIItem         *http.API
	ServerHTTPOrder *srvhttp.Config[*orderhttp.API] `ecfg:"SERVER_HTTP_ORDER"`
	APIOrder        *orderhttp.API
}

func isImplements[T any](v any) bool {
	_, ok := v.(T)
	return ok
}

func validateTypeResourcerOrConfigurable(e res.Entry) error {
	isResourcer := isImplements[Resourcer](e.Value)
	isConfigurable := isImplements[Configurable](e.Value)

	if isResourcer && isConfigurable {
		return fmt.Errorf("fail")
	}
	if !isResourcer && !isConfigurable {
		return fmt.Errorf("fail")
	}
	return nil
}

var (
	Pipeline        func()
	defaultPipeline = func() {
		// instance of global unique registry which contains added on init phase entries
		registryMain := unique.Global()
		// enforcement that registryMain contains Resourcer or Configurable type only
		registryMain.WalkEntries(func(e res.Entry) bool {
			if err := validateTypeResourcerOrConfigurable(e); err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			return true
		})

		// create new unique registry for ecfg work results
		registryResources := unique.New()
		registrySpecs := unique.New()
		// идем по appResources
		// если ConfigBuilder то создаем Spec и добавляем в registrySpecs и добавляем в registrySpecs user tag ecfg.TagECFG
		// если нет, то проверяем что это Resourcer и добавляем в registryResources

		// таким образом екфг работает с уник регситри а внем одлжны уже лежать просто спеки

		err := ecfg.Apply(registrySpecs)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		registrySpecs.WalkEntries(func(e res.Entry) bool {
			if !isImplements[Spec](e.Value) {
				fmt.Println(err)
				os.Exit(1)
			}

			built, err := e.Value.(Spec).Build()
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			if err := registryResources.Add(built); err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			return true
		})

		unique.Merge(registryMain, registryResources)

		// done
	}
)

func main() {
	if Pipeline != nil {
		defaultPipeline = Pipeline
	}
	defaultPipeline()
}
