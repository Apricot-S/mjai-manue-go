package mjairuntime

import (
	"bufio"
	"fmt"
	"io"

	"github.com/Apricot-S/mjai-manue-go/internal/adapter/mjai/inbound"
	"github.com/Apricot-S/mjai-manue-go/internal/adapter/mjai/outbound"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/ai"
)

type jsonLinesPolicy struct {
	respondNoneOnNoReaction bool
	stopOnEndGame           bool
	errorOnEOFBeforeEndGame bool
}

func runJSONLines(name string, room string, agent ai.Agent, in io.Reader, out io.Writer, log io.Writer, policy jsonLinesPolicy) error {
	r := bufio.NewScanner(in)
	w := bufio.NewWriter(out)
	defer w.Flush()

	driver := NewDriver(name, room, agent, log)
	for r.Scan() {
		stop, err := handleJSONLine(r.Bytes(), w, driver, log, policy)
		if err != nil {
			return err
		}
		if stop {
			return nil
		}
	}
	if err := r.Err(); err != nil {
		return err
	}
	if policy.errorOnEOFBeforeEndGame && !driver.Ended() {
		return fmt.Errorf("connection closed before end_game")
	}
	return nil
}

func handleJSONLine(line []byte, w *bufio.Writer, driver *Driver, log io.Writer, policy jsonLinesPolicy) (bool, error) {
	if err := traceLine(log, "<-", line); err != nil {
		return false, err
	}
	if len(line) == 0 {
		return false, fmt.Errorf("empty input line")
	}

	msg, err := inbound.ParseMessage(line)
	if err != nil {
		return false, err
	}
	outMsg, err := driver.Handle(msg)
	if err != nil {
		return false, err
	}
	if driver.Ended() && policy.stopOnEndGame {
		return true, nil
	}
	if outMsg == nil {
		if !policy.respondNoneOnNoReaction {
			return false, nil
		}
		outMsg = outbound.NewNone()
	}
	return false, writeMessageWithTrace(w, outMsg, log)
}
