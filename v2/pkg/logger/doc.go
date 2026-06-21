/*
Package logger provides structured logging (log/slog) with resources in res and SDI wiring.

# Bootstrap

Blank-import the use subpackage at the app composition root (not logger itself):

	import _ "github.com/omcrgnt/demo/v2/pkg/logger/use"

	logger.Info(ctx, "started", "port", 8080)

logger/use registers default configs in res ([res.TagReplaceable]) via [DefaultLogConfig] and [DefaultStdoutConfig].
Without logger/use and without [Config] in app config, package funcs are no-op until [sdi.Resolve] after [github.com/omcrgnt/demo/v2/pkg/builder].Build.

# User override

	type AppConfig struct {
	    Logger    logger.Config              `ecfg:"LOGGER"`
	    LogStdout logger.OutputStdoutConfig  `ecfg:"LOG_STDOUT"`
	    LogFile   logger.OutputFileConfig    `ecfg:"LOG_FILE"`
	}

Pipeline: ecfg.Register(cfg, res.Default) → builder.Build(res.Default) → sdi.Resolve(res.Default).
Dedup removes system defaults when user Logger or Output is registered.

# Usage

Package funcs (Info, Debug, …) and [Default] delegate to the logger wired by sdi (Inject sets active).

For DI, declare the port in Deps:

	func (c *Controller) Deps() []any {
	    return []any{(*logger.Logger)(nil)}
	}

If [Logger] is not in res, sdi.Resolve fails for resources that depend on it.

# SDI

Concrete logger implements [Logger]; Deps returns (*Output)(nil).
[Config.Build], [OutputStdoutConfig.Build], and [OutputFileConfig.Build] materialize runtime resources.
[DefaultLogConfig] and [DefaultStdoutConfig] are for logger/use system registration.
[DefaultLog] and [DefaultStdout] remain for tests.

See https://github.com/omcrgnt/demo/blob/main/docs/res-sdi-coupling.md.
*/
package logger
