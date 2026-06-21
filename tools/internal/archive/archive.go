package archive

import (
	"bufio"
	"bytes"
	"fmt"

	"github.com/Apricot-S/mjai-manue-go/internal/adapter/mjai/inbound"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/common"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/event"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/round"
)

const InitialScore = 25000

type Handlers struct {
	OnRaw      func(line []byte) error
	OnMessage  func(msg inbound.Message) error
	OnEvent    func(ev event.Event, archive *Archive) error
	OnFileDone func(path string) error
}

type Archive struct {
	state  *round.State
	scores [common.NumPlayers]int
}

func NewArchive() *Archive {
	a := &Archive{}
	a.resetScores()
	return a
}

func (a *Archive) State() (*round.State, bool) {
	if a.state == nil {
		return nil, false
	}
	return a.state, true
}

func (a *Archive) StateViewer() (round.StateViewer, bool) {
	return a.State()
}

func (a *Archive) Scores() [common.NumPlayers]int {
	return a.scores
}

func (a *Archive) PlayPaths(paths []string, h Handlers) error {
	for _, p := range paths {
		if err := a.playFile(p, h); err != nil {
			return err
		}
		if h.OnFileDone != nil {
			if err := h.OnFileDone(p); err != nil {
				return fmt.Errorf("%s: file-done callback failed: %w", p, err)
			}
		}
	}
	return nil
}

func (a *Archive) playFile(path string, h Handlers) error {
	reader, err := openMaybeGzip(path)
	if err != nil {
		return err
	}
	defer reader.Close()

	scanner := bufio.NewScanner(reader)
	lineNumber := 0
	for scanner.Scan() {
		lineNumber++
		line := bytes.TrimSpace(scanner.Bytes())
		if len(line) == 0 {
			return fmt.Errorf("%s:%d: empty line", path, lineNumber)
		}
		if err := a.processLine(line, h); err != nil {
			return fmt.Errorf("%s:%d: %w", path, lineNumber, err)
		}
	}
	if err := scanner.Err(); err != nil {
		return fmt.Errorf("%s: scanner error: %w", path, err)
	}
	return nil
}

func (a *Archive) processLine(line []byte, h Handlers) error {
	if h.OnRaw != nil {
		if err := h.OnRaw(bytes.Clone(line)); err != nil {
			return fmt.Errorf("raw callback failed: %w", err)
		}
	}

	msg, err := inbound.ParseMessage(line)
	if err != nil {
		return fmt.Errorf("failed to parse message: %w", err)
	}

	a.handleLifecycleMessage(msg)

	if h.OnMessage != nil {
		if err := h.OnMessage(msg); err != nil {
			return fmt.Errorf("message callback failed: %w", err)
		}
	}

	ev, err := inbound.ParseEvent(msg)
	if err != nil {
		return nil
	}

	if err := a.applyEvent(ev); err != nil {
		return err
	}

	if h.OnEvent != nil {
		if err := h.OnEvent(ev, a); err != nil {
			return fmt.Errorf("event callback failed: %w", err)
		}
	}

	if _, ok := ev.(*event.EndRound); ok {
		a.state = nil
	}
	return nil
}

func (a *Archive) handleLifecycleMessage(msg inbound.Message) {
	switch msg.(type) {
	case *inbound.StartGame:
		a.resetScores()
		a.state = nil
	case *inbound.EndGame:
		a.state = nil
	}
}

func (a *Archive) applyEvent(ev event.Event) error {
	switch ev := ev.(type) {
	case *event.StartRound:
		state, err := round.NewState(ev, a.scores)
		if err != nil {
			return fmt.Errorf("failed to start round: %w", err)
		}
		a.state = state
		a.scores = state.Scores()
		return nil
	case *event.EndRound:
		return nil
	default:
		if a.state == nil {
			return fmt.Errorf("cannot apply %T before start_kyoku", ev)
		}
		if err := a.state.Apply(ev); err != nil {
			return fmt.Errorf("failed to apply event: %w", err)
		}
		a.scores = a.state.Scores()
		return nil
	}
}

func (a *Archive) resetScores() {
	for i := range a.scores {
		a.scores[i] = InitialScore
	}
}
