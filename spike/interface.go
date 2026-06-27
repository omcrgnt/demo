package main

type Configurable interface {
	BuildConfig() (Spec, error)
}

type Spec interface {
	Build() (any, error)
}

type Resourcer interface {
	NewResource() (any, error)
}
