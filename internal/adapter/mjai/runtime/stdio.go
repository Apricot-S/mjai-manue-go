package mjairuntime

import (
	"io"

	"github.com/Apricot-S/mjai-manue-go/internal/domain/ai"
)

type StdioConfig struct {
	Name       string
	Room       string
	FallbackID int
	Agent      ai.Agent
	In         io.Reader
	Out        io.Writer
	Log        io.Writer
}

func RunStdio(cfg StdioConfig) error {
	return runJSONLines(cfg.Name, cfg.Room, cfg.FallbackID, cfg.Agent, cfg.In, cfg.Out, cfg.Log, jsonLinesPolicy{})
}
