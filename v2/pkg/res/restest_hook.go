package res

// ResetGlobalForRestest replaces [Global] with an empty registry.
//
//go:internal github.com/omcrgnt/demo/v2/pkg/res/restest
func ResetGlobalForRestest() Registry {
	resetGlobalRegistry()
	return global
}
