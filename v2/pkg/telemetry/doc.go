/*
Package telemetry provides OpenTelemetry trace export (OTLP) with resources in res and SDI wiring.

# Bootstrap

Blank-import the use subpackage at the app composition root (not telemetry itself):

	import _ "github.com/omcrgnt/demo/v2/pkg/telemetry/use"

telemetry/use registers a default [Config] in res ([res.TagReplaceable]) via [DefaultTraceConfig].
Hardcoded dev defaults: host localhost, port 4318, insecure TLS, service name app.

# User override

	type AppConfig struct {
	    Telemetry telemetry.Config `ecfg:"TELEMETRY"`
	}

[Config] uses [common.v1.Label], [common.v1.Host], and [common.v1.Port] for OTLP endpoint and service name.

Pipeline: ecfg.Register(cfg, res.Default) → builder.Build(res.Default) → sdi.Resolve(res.Default).
Dedup removes the system default when user Telemetry is registered.

# SDI and runner

[Provider] implements Inject (sets global TracerProvider) and Close (shutdown).
Resolve runs before runner.Run so HTTP otelhttp instrumentation sees the provider.

See https://github.com/omcrgnt/demo/blob/main/docs/res-sdi-coupling.md.
*/
package telemetry
