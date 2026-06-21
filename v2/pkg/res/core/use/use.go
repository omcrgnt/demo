// Package use registers org system defaults (logger, telemetry) in res.Default.
//
// Import for side effects at the app composition root:
//
//	import _ "github.com/omcrgnt/demo/v2/pkg/res/core/use"
package use

import (
	_ "github.com/omcrgnt/demo/v2/pkg/logger/use"
	_ "github.com/omcrgnt/demo/v2/pkg/telemetry/use"
)
