package main

import (
	"bufio"
	"bytes"
	"encoding/json/v2"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/Apricot-S/mjai-manue-go/configs"
	"github.com/Apricot-S/mjai-manue-go/internal/adapter/mjai/inbound"
	"github.com/Apricot-S/mjai-manue-go/internal/adapter/mjai/outbound"
	"github.com/Apricot-S/mjai-manue-go/internal/application"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/ai"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/round"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/seat"
)

const (
	defaultPlayerName = "Manue014"
	defaultSeed       = uint64(0)

	exitOK       = 0
	exitMismatch = 1
	exitRunError = 2
)

type config struct {
	playerName string
	seed       uint64
	limit      int
	showMatch  bool
	patterns   []string
}

type summary struct {
	files      int
	decisions  int
	matches    int
	mismatches int
	errors     int
}

type normalizedAction struct {
	Type      string   `json:"type"`
	Actor     *int     `json:"actor,omitempty"`
	Target    *int     `json:"target,omitempty"`
	Pai       string   `json:"pai,omitempty"`
	Consumed  []string `json:"consumed,omitempty"`
	Tsumogiri *bool    `json:"tsumogiri,omitempty"`
	Reason    string   `json:"reason,omitempty"`
}

func main() {
	os.Exit(run(os.Args[1:], os.Stdout, os.Stderr))
}

func run(args []string, out io.Writer, errOut io.Writer) int {
	cfg, err := parseConfig(args, errOut)
	if err != nil {
		fmt.Fprintln(errOut, err)
		return exitRunError
	}

	paths, err := globAll(cfg.patterns)
	if err != nil {
		fmt.Fprintln(errOut, err)
		return exitRunError
	}

	stats, err := configs.LoadGameStats()
	if err != nil {
		fmt.Fprintf(errOut, "failed to load game stats: %v\n", err)
		return exitRunError
	}
	dangerTree, err := configs.LoadDangerTree()
	if err != nil {
		fmt.Fprintf(errOut, "failed to load danger tree: %v\n", err)
		return exitRunError
	}

	c := comparer{
		cfg: cfg,
		deps: ai.ManueAgentDeps{
			Stats:  stats,
			Danger: ai.NewDangerEstimator(dangerTree),
		},
		out: out,
		log: errOut,
	}

	s := summary{files: len(paths)}
	for _, path := range paths {
		fileSummary, err := c.compareFile(path)
		s.add(fileSummary)
		if err != nil {
			s.errors++
			fmt.Fprintf(out, "error: %s: %v\n", path, err)
		}
	}
	fmt.Fprintf(out, "summary: files=%d decisions=%d matches=%d mismatches=%d errors=%d\n",
		s.files, s.decisions, s.matches, s.mismatches, s.errors)
	if s.errors > 0 {
		return exitRunError
	}
	if s.mismatches > 0 {
		return exitMismatch
	}
	return exitOK
}

func parseConfig(args []string, errOut io.Writer) (config, error) {
	flags := flag.NewFlagSet("compare", flag.ContinueOnError)
	flags.SetOutput(errOut)
	playerName := flags.String("player-name", defaultPlayerName, "original player name")
	seed := flags.Uint64("seed", defaultSeed, "Go port random seed")
	limit := flags.Int("limit", 0, "maximum number of mismatches to report; 0 means unlimited")
	showMatch := flags.Bool("show-matches", false, "print matched decisions")
	if err := flags.Parse(args); err != nil {
		return config{}, err
	}
	if flags.NArg() == 0 {
		return config{}, errors.New("usage: compare [OPTIONS] <LOG_GLOB_PATTERNS>...")
	}
	if *limit < 0 {
		return config{}, errors.New("--limit must be >= 0")
	}
	return config{
		playerName: *playerName,
		seed:       *seed,
		limit:      *limit,
		showMatch:  *showMatch,
		patterns:   flags.Args(),
	}, nil
}

func globAll(patterns []string) ([]string, error) {
	var paths []string
	for _, pattern := range patterns {
		matches, err := filepath.Glob(pattern)
		if err != nil {
			return nil, fmt.Errorf("invalid glob pattern %q: %w", pattern, err)
		}
		if len(matches) == 0 {
			if _, err := os.Stat(pattern); err == nil {
				matches = []string{pattern}
			}
		}
		if len(matches) == 0 {
			return nil, fmt.Errorf("no files match %q", pattern)
		}
		paths = append(paths, matches...)
	}
	slices.Sort(paths)
	return slices.Compact(paths), nil
}

func (s *summary) add(other summary) {
	s.decisions += other.decisions
	s.matches += other.matches
	s.mismatches += other.mismatches
	s.errors += other.errors
}

type comparer struct {
	cfg  config
	deps ai.ManueAgentDeps
	out  io.Writer
	log  io.Writer
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
		if reaction.Kind() == application.ReactionAction {
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
		}
	}
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
		fmt.Fprintf(fc.parent.out, "original: %s\n", mustActionJSON(original))
	} else {
		fmt.Fprintln(fc.parent.out, "original: <none>")
	}
	if port != nil {
		fmt.Fprintf(fc.parent.out, "port:     %s\n", mustActionJSON(*port))
	} else {
		fmt.Fprintln(fc.parent.out, "port:     <none>")
	}
	if fc.lastLog != "" {
		fmt.Fprintf(fc.parent.out, "original log:\n%s\n", fc.lastLog)
	}
	if fc.lastBoard != "" {
		fmt.Fprintf(fc.parent.out, "board:\n%s", fc.lastBoard)
	}
	if fc.lastTrace != "" {
		fmt.Fprintf(fc.parent.out, "trace:\n%s", fc.lastTrace)
	}
	fmt.Fprintln(fc.parent.out, strings.Repeat("-", 80))
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

func normalizeRawAction(raw []byte) (normalizedAction, bool, error) {
	var m map[string]any
	if err := json.Unmarshal(raw, &m); err != nil {
		return normalizedAction{}, false, err
	}
	t, _ := m["type"].(string)
	action := normalizedAction{Type: t}

	switch t {
	case "dahai":
		action.Actor = intPtrFromMap(m, "actor")
		action.Pai, _ = m["pai"].(string)
		action.Tsumogiri = boolPtrFromMap(m, "tsumogiri")
	case "reach":
		action.Actor = intPtrFromMap(m, "actor")
	case "chi", "pon", "daiminkan":
		action.Actor = intPtrFromMap(m, "actor")
		action.Target = intPtrFromMap(m, "target")
		action.Pai, _ = m["pai"].(string)
		action.Consumed = stringSliceFromMap(m, "consumed")
	case "ankan":
		action.Actor = intPtrFromMap(m, "actor")
		action.Consumed = stringSliceFromMap(m, "consumed")
	case "kakan":
		action.Actor = intPtrFromMap(m, "actor")
		action.Pai, _ = m["pai"].(string)
		action.Consumed = stringSliceFromMap(m, "consumed")
	case "hora":
		action.Actor = intPtrFromMap(m, "actor")
		action.Target = intPtrFromMap(m, "target")
		action.Pai, _ = m["pai"].(string)
	case "ryukyoku":
		action.Actor = intPtrFromMap(m, "actor")
		action.Reason, _ = m["reason"].(string)
		if action.Reason != "kyushukyuhai" {
			return normalizedAction{}, false, nil
		}
	case "none":
		action.Actor = intPtrFromMap(m, "actor")
	default:
		return normalizedAction{}, false, nil
	}
	return action, true, nil
}

func intPtrFromMap(m map[string]any, key string) *int {
	v, ok := m[key]
	if !ok {
		return nil
	}
	switch n := v.(type) {
	case int:
		return &n
	case int64:
		i := int(n)
		return &i
	case float64:
		i := int(n)
		return &i
	}
	return nil
}

func boolPtrFromMap(m map[string]any, key string) *bool {
	v, ok := m[key].(bool)
	if !ok {
		return nil
	}
	return &v
}

func stringSliceFromMap(m map[string]any, key string) []string {
	values, ok := m[key].([]any)
	if !ok {
		return nil
	}
	out := make([]string, 0, len(values))
	for _, value := range values {
		s, ok := value.(string)
		if !ok {
			continue
		}
		out = append(out, s)
	}
	return out
}

func actionsEqual(a, b normalizedAction) bool {
	if a.Type != b.Type || !intPtrEqual(a.Actor, b.Actor) || !intPtrEqual(a.Target, b.Target) {
		return false
	}
	if a.Pai != b.Pai || !boolPtrEqual(a.Tsumogiri, b.Tsumogiri) || a.Reason != b.Reason {
		return false
	}
	return slices.Equal(a.Consumed, b.Consumed)
}

func intPtrEqual(a, b *int) bool {
	if a == nil || b == nil {
		return a == b
	}
	return *a == *b
}

func boolPtrEqual(a, b *bool) bool {
	if a == nil || b == nil {
		return a == b
	}
	return *a == *b
}

func mustActionJSON(a normalizedAction) string {
	b, err := json.Marshal(a)
	if err != nil {
		return fmt.Sprintf("%+v", a)
	}
	return string(b)
}

func isZeroAction(a normalizedAction) bool {
	return a.Type == "" && a.Actor == nil && a.Target == nil && a.Pai == "" && len(a.Consumed) == 0 && a.Tsumogiri == nil && a.Reason == ""
}
