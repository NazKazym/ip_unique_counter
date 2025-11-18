package main

import (
	"context"
	"fmt"
	"io"
	"os"
)

// InputSource is an abstraction for where the data comes from.
// Right now we only implement local files, but more can be added later.
type InputSource interface {
	Open(ctx context.Context) (io.ReadCloser, error)
}

// FileSource implements InputSource for a local filesystem path.
type FileSource struct {
	Path string
}

func NewFileSource(path string) *FileSource {
	return &FileSource{Path: path}
}

func (f *FileSource) Open(ctx context.Context) (io.ReadCloser, error) {
	// ctx is available if you want cancellation logic later.
	return os.Open(f.Path)
}

// BuildInputSource selects the correct InputSource implementation based on config.
func BuildInputSource(cfg Config) (InputSource, error) {
	switch cfg.SourceType {
	case "", defaultSourceTypeFile:
		return NewFileSource(cfg.SourceURI), nil
	default:
		return nil, fmt.Errorf("unsupported source_type: %s (only 'file' is implemented)", cfg.SourceType)
	}
}
