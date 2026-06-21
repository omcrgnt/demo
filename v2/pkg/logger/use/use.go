// Package use registers logger system defaults in res.Default.
//
// Import for side effects at the app composition root (main or a meta use package):
//
//	import _ "github.com/omcrgnt/demo/v2/pkg/logger/use"
package use

import (
	"github.com/omcrgnt/demo/v2/pkg/logger"
	"github.com/omcrgnt/demo/v2/pkg/res"
)

func init() {
	_ = res.AddWithTags(logger.DefaultLogConfig(), res.TagReplaceable)
	_ = res.AddWithTags(logger.DefaultStdoutConfig(), res.TagReplaceable)
}
