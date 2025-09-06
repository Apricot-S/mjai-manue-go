package shared

import (
	"bufio"
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/Apricot-S/mjai-manue-go/internal/game"
	"github.com/Apricot-S/mjai-manue-go/internal/game/event/inbound"
	"github.com/Apricot-S/mjai-manue-go/internal/protocol"
	"github.com/schollz/progressbar/v3"
)

type Archive struct {
	paths   []string
	adapter protocol.Adapter
	state   game.State
}

func NewArchive(paths []string, adapter protocol.Adapter) *Archive {
	return &Archive{
		paths:   paths,
		adapter: adapter,
		state:   &game.StateImpl{},
	}
}

func (a *Archive) StateViewer() game.StateViewer {
	return a.state
}

func (a *Archive) StateUpdater() game.StateUpdater {
	return a.state
}

func (a *Archive) StateAnalyzer() game.StateAnalyzer {
	return a.state
}

func (a *Archive) PlayLight(onAction func(inbound.Event) error) error {
	numFiles := len(a.paths)
	bar := progressbar.Default(int64(numFiles))
	for _, p := range a.paths {
		if err := a.playLightInner(p, onAction); err != nil {
			return err
		}

		if err := bar.Add(1); err != nil {
			return err
		}
	}
	return nil
}

func (a *Archive) playLightInner(singlePath string, onAction func(inbound.Event) error) error {
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

		action, err := a.parseAction(line)
		if err != nil {
			return err
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

func (a *Archive) parseAction(line []byte) (inbound.Event, error) {
	actions, err := a.adapter.DecodeMessages(line)
	if err != nil {
		return nil, fmt.Errorf("failed to decode: %w", err)
	}
	if len(actions) != 1 {
		return nil, fmt.Errorf("expected 1 action in 1 line, got %d", len(actions))
	}
	return actions[0], nil
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

func (a *Archive) Play(onAction func(inbound.Event) error) error {
	onLightAction := func(action inbound.Event) error {
		if err := a.state.Update(action); err != nil {
			return err
		}
		return onAction(action)
	}
	return a.PlayLight(onLightAction)
}
