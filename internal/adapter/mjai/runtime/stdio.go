package mjairuntime

import (
	"bufio"
	"fmt"
	"io"

	"github.com/Apricot-S/mjai-manue-go/internal/adapter/mjai/inbound"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/ai"
)

type StdioConfig struct {
	Name  string
	Room  string
	Agent ai.Agent
	In    io.Reader
	Out   io.Writer
	Log   io.Writer
}

func RunStdio(cfg StdioConfig) error {
	r := bufio.NewScanner(cfg.In)
	w := bufio.NewWriter(cfg.Out)
	defer w.Flush()

	driver := NewDriver(cfg.Name, cfg.Room, cfg.Agent, cfg.Log)
	for r.Scan() {
		line := r.Bytes()
		if err := traceLine(cfg.Log, "<-", line); err != nil {
			return err
		}
		if len(line) == 0 {
			return fmt.Errorf("empty input line")
		}

		msg, err := inbound.ParseMessage(line)
		if err != nil {
			return err
		}
		outMsg, err := driver.Handle(msg)
		if err != nil {
			return err
		}
		if outMsg == nil {
			continue
		}
		if err := writeMessageWithTrace(w, outMsg, cfg.Log); err != nil {
			return err
		}
	}
	if err := r.Err(); err != nil {
		return err
	}
	return nil
}
