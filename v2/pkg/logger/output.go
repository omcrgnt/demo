package logger

import "io"

// Output is a sealed log sink port for SDI.
type Output interface {
	io.Writer
	markerLoggerOutput()
}
