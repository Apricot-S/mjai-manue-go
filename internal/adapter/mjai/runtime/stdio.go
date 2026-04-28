package mjairuntime

import (
	"bufio"
	"fmt"
	"io"

	"github.com/Apricot-S/mjai-manue-go/internal/adapter/mjai/inbound"
	"github.com/Apricot-S/mjai-manue-go/internal/adapter/mjai/outbound"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/ai"
)

type StdioConfig struct {
	Name  string
	Room  string
	Agent ai.Agent
	In    io.Reader
	Out   io.Writer
}

func RunStdio(cfg StdioConfig) error {
	r := bufio.NewScanner(cfg.In)
	w := bufio.NewWriter(cfg.Out)
	defer w.Flush()

	driver := NewDriver(cfg.Name, cfg.Room, cfg.Agent)
	for r.Scan() {
		line := r.Bytes()
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
		if err := writeMessage(w, outMsg); err != nil {
			return err
		}
	}
	if err := r.Err(); err != nil {
		return err
	}
	return nil
}

func writeMessage(w *bufio.Writer, msg outbound.Message) error {
	b, err := outbound.MarshalMessage(msg)
	if err != nil {
		return err
	}
	if _, err := w.Write(b); err != nil {
		return err
	}
	if err := w.WriteByte('\n'); err != nil {
		return err
	}
	return w.Flush()
}
