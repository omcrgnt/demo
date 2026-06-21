package logger

import (
	"fmt"
	"os"
)

// OutputFileConfig opens a file log sink (user resource via Build).
type OutputFileConfig struct {
	Path string `ecfg:"PATH"`
}

type fileOutput struct {
	f *os.File
}

func (f *fileOutput) Write(p []byte) (int, error) {
	if f.f == nil {
		return 0, fmt.Errorf("logger: file output closed")
	}
	return f.f.Write(p)
}

func (f *fileOutput) markerLoggerOutput() {}

// Build returns an Output writing to Path.
func (c OutputFileConfig) Build() (any, error) {
	if c.Path == "" {
		return nil, fmt.Errorf("logger: output file path is required")
	}
	f, err := os.OpenFile(c.Path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o644)
	if err != nil {
		return nil, fmt.Errorf("logger: open output file: %w", err)
	}
	return &fileOutput{f: f}, nil
}
