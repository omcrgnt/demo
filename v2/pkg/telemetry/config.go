package telemetry

import (
	commonv1 "github.com/omcrgnt/proto/gen/go/common/v1"
)

// Config configures OTLP trace export (host, port, service name).
// ecfg fills and protovalidates proto fields before Build is called.
type Config struct {
	ServiceName *commonv1.Label `ecfg:"SERVICE_NAME"`
	Host        *commonv1.Host  `ecfg:"HOST"`
	Port        *commonv1.Port  `ecfg:"PORT"`
	Insecure    bool            `ecfg:"INSECURE"`
}

// Build returns a Provider for res.Add and SDI.
func (c Config) Build() (any, error) {
	return newProvider(
		c.Host.GetValue(),
		c.Port.GetValue(),
		c.Insecure,
		c.ServiceName.GetValue(),
	)
}
