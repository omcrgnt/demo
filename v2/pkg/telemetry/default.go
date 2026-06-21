package telemetry

import commonv1 "github.com/omcrgnt/proto/gen/go/common/v1"

const (
	defaultHostValue        = "localhost"
	defaultPortValue        = uint32(4318)
	defaultServiceNameValue = "app"
	defaultInsecureValue    = true
)

// DefaultTraceConfig returns the system trace config for telemetry/use registration.
func DefaultTraceConfig() Config {
	return Config{
		ServiceName: &commonv1.Label{Value: defaultServiceNameValue},
		Host:        &commonv1.Host{Value: defaultHostValue},
		Port:        &commonv1.Port{Value: defaultPortValue},
		Insecure:    defaultInsecureValue,
	}
}

// DefaultTrace returns the system trace Provider for tests and legacy callers.
func DefaultTrace() any {
	p, err := newProvider(defaultHostValue, defaultPortValue, defaultInsecureValue, defaultServiceNameValue)
	if err != nil {
		return &Provider{tp: noopTracerProvider()}
	}
	return p
}
