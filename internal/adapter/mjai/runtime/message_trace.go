package mjairuntime

import (
	"bufio"
	"fmt"
	"io"

	"github.com/Apricot-S/mjai-manue-go/internal/adapter/mjai/outbound"
)

func writeMessageWithTrace(w *bufio.Writer, msg outbound.Message, log io.Writer) error {
	b, err := outbound.MarshalMessage(msg)
	if err != nil {
		return err
	}
	if err := traceLine(log, "->", b); err != nil {
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

func traceLine(log io.Writer, direction string, line []byte) error {
	if log == nil {
		return nil
	}
	_, err := fmt.Fprintf(log, "%s\t%s\n", direction, string(line))
	return err
}
