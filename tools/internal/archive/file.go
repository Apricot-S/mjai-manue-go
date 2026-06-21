package archive

import (
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

func openMaybeGzip(path string) (io.ReadCloser, error) {
	file, err := os.Open(filepath.Clean(path))
	if err != nil {
		return nil, fmt.Errorf("failed to open %s: %w", path, err)
	}

	if filepath.Ext(path) != ".gz" {
		return file, nil
	}

	gz, err := gzip.NewReader(file)
	if err != nil {
		file.Close()
		return nil, fmt.Errorf("failed to open gzip %s: %w", path, err)
	}

	return &gzipReadCloser{Reader: gz, file: file}, nil
}

type gzipReadCloser struct {
	*gzip.Reader
	file *os.File
}

func (g *gzipReadCloser) Close() error {
	err1 := g.Reader.Close()
	err2 := g.file.Close()
	if err1 != nil {
		return err1
	}
	return err2
}
