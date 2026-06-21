package obs

// ApplyTransform wraps v with Observe() when v implements [Observer];
// otherwise returns v unchanged. Use with res.Transform:
//
//	res.Transform(obs.ApplyTransform)
func ApplyTransform(v any) any {
	if o, ok := v.(Observer); ok {
		return o.Observe()
	}
	return v
}
