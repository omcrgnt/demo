// Package use registers telemetry system defaults in res.Global().
//
// Import for side effects at the app composition root:
//
//	import _ "github.com/omcrgnt/demo/v2/pkg/telemetry/use"
package use

import (
	"github.com/omcrgnt/demo/v2/pkg/res"
	"github.com/omcrgnt/demo/v2/pkg/telemetry"
)

func init() {
	_ = res.AddToGlobalWithTags(telemetry.DefaultTraceConfig(), res.TagReplaceable)
}
