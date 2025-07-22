package shared

import (
	"bufio"
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/go-json-experiment/json"
	"github.com/go-json-experiment/json/jsontext"
)

type Archive struct {
	paths []string
}

func NewArchive(paths []string) *Archive {
	return &Archive{paths: paths}
}

func (a *Archive) PlayLight(onAction func(jsontext.Value) error) error {
	numFiles := len(a.paths)
	for i, p := range a.paths {
		if numFiles > 1 {
			fmt.Fprintf(os.Stderr, "%d/%d\n", i+1, numFiles) // tentative
		}

		if err := playLightInner(p, onAction); err != nil {
			return err
		}
	}
	return nil
}

func playLightInner(singlePath string, onAction func(jsontext.Value) error) error {
	reader, err := openMaybeGzip(singlePath)
	if err != nil {
		return err
	}
	defer reader.Close()

	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		line := bytes.TrimSpace(scanner.Bytes())
		if len(line) == 0 {
			continue
		}

		var action jsontext.Value
		if err := json.Unmarshal(line, &action); err != nil {
			return fmt.Errorf("json decode error: %w", err)
		}

		if err := onAction(action); err != nil {
			return fmt.Errorf("failed to callback: %w", err)
		}
	}
	if err := scanner.Err(); err != nil {
		return fmt.Errorf("scanner error: %w", err)
	}
	return nil
}

func openMaybeGzip(path string) (io.ReadCloser, error) {
	file, err := os.Open(filepath.Clean(path))
	if err != nil {
		return nil, fmt.Errorf("failed to open %s: %w", path, err)
	}

	if filepath.Ext(path) != ".gz" {
		return file, nil // Return as a regular file
	}

	gz, err := gzip.NewReader(file)
	if err != nil {
		file.Close()
		return nil, fmt.Errorf("gzip error: %w", err)
	}

	return &gzipReadCloser{gz, file}, nil
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

// GetLightActions collects all actions from files into a slice.
func (a *Archive) GetLightActions() ([]jsontext.Value, error) {
	var actions []jsontext.Value
	err := a.PlayLight(func(action jsontext.Value) error {
		actions = append(actions, action)
		return nil
	})
	return actions, err
}
