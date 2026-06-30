package main

import (
	"bufio"
	"bytes"
	"encoding/json/v2"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/Apricot-S/mjai-manue-go/internal/adapter/mjai/inbound"
	"github.com/Apricot-S/mjai-manue-go/internal/adapter/mjai/outbound"
	"github.com/Apricot-S/mjai-manue-go/internal/application"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/ai"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/round"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/seat"
)

type comparer struct {
	cfg  config
	deps ai.ManueAgentDeps
	out  io.Writer
	log  io.Writer
}

const mismatchSeparatorWidth = 122

func (c *comparer) compareFile(path string) (summary, error) {
	agent, err := ai.NewManueAgent(c.cfg.seed, c.deps)
	if err != nil {
		return summary{}, fmt.Errorf("failed to create ManueAgent: %w", err)
	}

	fc := &fileComparer{
		parent: c,
		path:   path,
		self:   -1,
	}

	file, err := os.Open(path)
	if err != nil {
		return summary{}, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Buffer(make([]byte, 64*1024), 16*1024*1024)
	lineNo := 0
	for scanner.Scan() {
		lineNo++
		raw := bytes.TrimSpace(scanner.Bytes())
		if len(raw) == 0 {
			return fc.fileSummary, fmt.Errorf("line %d: empty line", lineNo)
		}
		if err := fc.processLine(lineNo, raw, agent); err != nil {
			return fc.fileSummary, err
		}
	}
	if err := scanner.Err(); err != nil {
		return fc.fileSummary, err
	}
	if err := fc.flushPendingAtEOF(); err != nil {
		return fc.fileSummary, err
	}
	return fc.fileSummary, nil
}

type fileComparer struct {
	parent      *comparer
	path        string
	bot         *application.Bot
	pending     *pendingAction
	self        int
	started     bool
	lastLog     string
	lastBoard   string
	lastTrace   string
	reported    int
	fileSummary summary
}

type pendingAction struct {
	line   int
	action normalizedAction
	raw    string
}

func (fc *fileComparer) processLine(lineNo int, raw []byte, agent ai.Agent) error {
	msg, err := inbound.ParseMessage(raw)
	if err != nil {
		return fmt.Errorf("line %d: parse message: %w", lineNo, err)
	}
	fc.captureOriginalLog(raw)

	original, originalComparable, err := normalizeRawAction(raw)
	if err != nil {
		return fmt.Errorf("line %d: normalize original action: %w", lineNo, err)
	}
	if originalComparable && original.Actor != nil && *original.Actor == fc.self {
		fc.compareOriginalSelfAction(lineNo, original)
	} else if err := fc.flushPendingBeforeNonSelf(lineNo, originalComparable); err != nil {
		return err
	}

	switch msg := msg.(type) {
	case *inbound.StartGame:
		self, err := findPlayer(msg.Names, fc.parent.cfg.playerName)
		if err != nil {
			return fmt.Errorf("line %d: %w", lineNo, err)
		}
		selfSeat, err := seat.NewSeat(self)
		if err != nil {
			return fmt.Errorf("line %d: invalid self seat: %w", lineNo, err)
		}
		agent.Reset()
		fc.self = self
		fc.started = true
		fc.bot = application.NewBot(selfSeat, agent, fc)
	case *inbound.EndGame:
		if err := fc.flushPendingAtEOF(); err != nil {
			return err
		}
		fc.started = false
		fc.bot = nil
		fc.self = -1
	default:
		if err := fc.processEvent(lineNo, msg); err != nil {
			return err
		}
	}
	return nil
}

func (fc *fileComparer) processEvent(lineNo int, msg inbound.Message) error {
	if !fc.started || fc.bot == nil {
		return fmt.Errorf("line %d: cannot process %T before start_game", lineNo, msg)
	}
	ev, err := inbound.ParseEvent(msg)
	if err != nil {
		return fmt.Errorf("line %d: parse event: %w", lineNo, err)
	}
	reaction, err := fc.bot.Process(ev)
	if err != nil {
		return fmt.Errorf("line %d: process event: %w", lineNo, err)
	}
	if reaction.Kind() != application.ReactionAction {
		return nil
	}

	actionMsg, err := outbound.ToMessage(reaction.Action(), reaction.Log())
	if err != nil {
		return fmt.Errorf("line %d: convert action: %w", lineNo, err)
	}
	actionRaw, err := outbound.MarshalMessage(actionMsg)
	if err != nil {
		return fmt.Errorf("line %d: marshal action: %w", lineNo, err)
	}
	action, _, err := normalizeRawAction(actionRaw)
	if err != nil {
		return fmt.Errorf("line %d: normalize Go action: %w", lineNo, err)
	}
	fc.pending = &pendingAction{line: lineNo, action: action, raw: string(actionRaw)}
	return nil
}

func findPlayer(names []string, playerName string) (int, error) {
	count := 0
	found := -1
	for i, name := range names {
		if name == playerName {
			count++
			found = i
		}
	}
	if count != 1 {
		return -1, fmt.Errorf("player name %q matched %d players in %v", playerName, count, names)
	}
	return found, nil
}

func (fc *fileComparer) compareOriginalSelfAction(lineNo int, original normalizedAction) {
	fc.fileSummary.decisions++
	if fc.pending == nil {
		fc.recordMismatch(lineNo, original, nil, "Go port returned no action")
		return
	}
	pending := fc.pending
	fc.pending = nil
	if actionsEqual(original, pending.action) {
		fc.fileSummary.matches++
		if fc.parent.cfg.showMatch {
			fmt.Fprintf(fc.parent.out, "match: %s:%d %s\n", fc.path, lineNo, mustActionJSON(original))
		}
		return
	}
	fc.recordMismatch(lineNo, original, &pending.action, "action mismatch")
}

func (fc *fileComparer) flushPendingBeforeNonSelf(lineNo int, comparable bool) error {
	if fc.pending == nil {
		return nil
	}
	pending := fc.pending
	fc.pending = nil
	if pending.action.Type == "none" {
		fc.fileSummary.decisions++
		fc.fileSummary.matches++
		if fc.parent.cfg.showMatch {
			fmt.Fprintf(fc.parent.out, "match: %s:%d implicit pass %s\n", fc.path, lineNo, mustActionJSON(pending.action))
		}
		return nil
	}
	if comparable {
		fc.fileSummary.decisions++
		fc.recordMismatch(lineNo, normalizedAction{}, &pending.action, "Go port returned an action, but original did not take it")
	}
	return nil
}

func (fc *fileComparer) flushPendingAtEOF() error {
	if fc.pending == nil {
		return nil
	}
	pending := fc.pending
	fc.pending = nil
	if pending.action.Type == "none" {
		fc.fileSummary.decisions++
		fc.fileSummary.matches++
		return nil
	}
	fc.fileSummary.decisions++
	fc.recordMismatch(pending.line, normalizedAction{}, &pending.action, "Go port returned an action at end of stream")
	return nil
}

func (fc *fileComparer) recordMismatch(lineNo int, original normalizedAction, port *normalizedAction, reason string) {
	fc.fileSummary.mismatches++
	if fc.parent.cfg.limit > 0 && fc.reported >= fc.parent.cfg.limit {
		return
	}
	fc.reported++

	fmt.Fprintf(fc.parent.out, "mismatch: %s:%d: %s\n", fc.path, lineNo, reason)
	if !isZeroAction(original) {
		fmt.Fprintf(fc.parent.out, "original action: %s\n", mustActionJSON(original))
	} else {
		fmt.Fprintln(fc.parent.out, "original action: <none>")
	}
	if port != nil {
		fmt.Fprintf(fc.parent.out, "port action:     %s\n", mustActionJSON(*port))
	} else {
		fmt.Fprintln(fc.parent.out, "port action:     <none>")
	}
	if fc.lastBoard != "" {
		fmt.Fprintf(fc.parent.out, "board:\n%s", fc.lastBoard)
	}
	if fc.lastLog != "" {
		fmt.Fprintf(fc.parent.out, "original log:\n%s\n", fc.lastLog)
	}
	if fc.lastTrace != "" {
		fmt.Fprintf(fc.parent.out, "port trace log:\n%s", fc.lastTrace)
	}
	fmt.Fprintln(fc.parent.out, strings.Repeat("=", mismatchSeparatorWidth))
}

func (fc *fileComparer) captureOriginalLog(raw []byte) {
	var payload struct {
		Logs []string `json:"logs"`
	}
	if err := json.Unmarshal(raw, &payload); err != nil {
		return
	}
	for _, log := range payload.Logs {
		if log != "" {
			fc.lastLog = log
		}
	}
}

func (fc *fileComparer) ReportRoundState(state round.BoardRenderer) error {
	fc.lastBoard = state.RenderBoard()
	return nil
}

func (fc *fileComparer) ReportDecisionTrace(trace string) error {
	fc.lastTrace = trace
	if trace != "" && fc.parent.log != nil {
		_, err := fmt.Fprint(fc.parent.log, trace)
		return err
	}
	return nil
}
