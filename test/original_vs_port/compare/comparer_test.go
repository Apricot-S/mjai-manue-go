package main

import (
	"bytes"
	"strings"
	"testing"
)

func TestFindPlayer(t *testing.T) {
	got, err := findPlayer([]string{"a", defaultPlayerName, "b", "c"}, defaultPlayerName)
	if err != nil {
		t.Fatalf("findPlayer() failed: %v", err)
	}
	if got != 1 {
		t.Errorf("player index = %d, want 1", got)
	}

	if _, err := findPlayer([]string{"a", "b", "c", "d"}, defaultPlayerName); err == nil {
		t.Error("findPlayer() with no match succeeded, want error")
	}
	if _, err := findPlayer([]string{defaultPlayerName, "b", defaultPlayerName, "d"}, defaultPlayerName); err == nil {
		t.Error("findPlayer() with duplicate matches succeeded, want error")
	}
}

func TestFileComparer_CompareOriginalSelfActionMatchesPending(t *testing.T) {
	var out bytes.Buffer
	fc := newTestFileComparer(&out, 0)
	fc.pending = &pendingAction{
		line:   10,
		action: normalizedAction{Type: "dahai", Actor: new(1), Pai: "5m", Tsumogiri: new(false)},
	}

	fc.compareOriginalSelfAction(11, normalizedAction{Type: "dahai", Actor: new(1), Pai: "5m", Tsumogiri: new(false)})

	if fc.pending != nil {
		t.Error("pending action was not cleared")
	}
	if fc.fileSummary.decisions != 1 || fc.fileSummary.matches != 1 || fc.fileSummary.implicitPasses != 0 || fc.fileSummary.mismatches != 0 {
		t.Errorf("summary = %+v, want 1 direct match only", fc.fileSummary)
	}
	if out.String() != "" {
		t.Errorf("output = %q, want empty", out.String())
	}
}

func TestFileComparer_FlushPendingNonNoneBeforeNonSelfReportsMismatch(t *testing.T) {
	var out bytes.Buffer
	fc := newTestFileComparer(&out, 0)
	fc.pending = &pendingAction{
		line:   10,
		action: normalizedAction{Type: "dahai", Actor: new(1), Pai: "5m", Tsumogiri: new(false)},
	}

	if err := fc.flushPendingBeforeNonSelf(11); err != nil {
		t.Fatalf("flushPendingBeforeNonSelf() failed: %v", err)
	}

	if fc.pending != nil {
		t.Error("pending action was not cleared")
	}
	if fc.fileSummary.decisions != 1 || fc.fileSummary.matches != 0 || fc.fileSummary.implicitPasses != 0 || fc.fileSummary.mismatches != 1 {
		t.Errorf("summary = %+v, want 1 mismatch only", fc.fileSummary)
	}
	if !strings.Contains(out.String(), "Go port returned an action, but original did not take it") {
		t.Errorf("output = %q, want mismatch reason", out.String())
	}
}

func TestFileComparer_FlushPendingNonNoneAtEOFReportsMismatch(t *testing.T) {
	var out bytes.Buffer
	fc := newTestFileComparer(&out, 0)
	fc.pending = &pendingAction{
		line:   10,
		action: normalizedAction{Type: "reach", Actor: new(1)},
	}

	if err := fc.flushPendingAtEOF(); err != nil {
		t.Fatalf("flushPendingAtEOF() failed: %v", err)
	}

	if fc.pending != nil {
		t.Error("pending action was not cleared")
	}
	if fc.fileSummary.decisions != 1 || fc.fileSummary.matches != 0 || fc.fileSummary.implicitPasses != 0 || fc.fileSummary.mismatches != 1 {
		t.Errorf("summary = %+v, want 1 mismatch only", fc.fileSummary)
	}
	if !strings.Contains(out.String(), "Go port returned an action at end of stream") {
		t.Errorf("output = %q, want EOF mismatch reason", out.String())
	}
}

func TestFileComparer_FlushPendingNoneCountsImplicitPass(t *testing.T) {
	var out bytes.Buffer
	fc := newTestFileComparer(&out, 0)
	fc.pending = &pendingAction{
		line:   10,
		action: normalizedAction{Type: "none", Actor: new(1)},
	}

	if err := fc.flushPendingBeforeNonSelf(11); err != nil {
		t.Fatalf("flushPendingBeforeNonSelf() failed: %v", err)
	}

	if fc.pending != nil {
		t.Error("pending action was not cleared")
	}
	if fc.fileSummary.decisions != 1 || fc.fileSummary.matches != 0 || fc.fileSummary.implicitPasses != 1 || fc.fileSummary.mismatches != 0 {
		t.Errorf("summary = %+v, want 1 implicit pass only", fc.fileSummary)
	}
}

func TestFileComparer_OriginalSelfActionWithoutPendingReportsMismatch(t *testing.T) {
	var out bytes.Buffer
	fc := newTestFileComparer(&out, 0)

	fc.compareOriginalSelfAction(11, normalizedAction{Type: "hora", Actor: new(1), Target: new(0), Pai: "5m"})

	if fc.fileSummary.decisions != 1 || fc.fileSummary.matches != 0 || fc.fileSummary.implicitPasses != 0 || fc.fileSummary.mismatches != 1 {
		t.Errorf("summary = %+v, want 1 mismatch only", fc.fileSummary)
	}
	if !strings.Contains(out.String(), "Go port returned no action") {
		t.Errorf("output = %q, want missing action reason", out.String())
	}
}

func TestFileComparer_MismatchLimitOnlyLimitsDetails(t *testing.T) {
	var out bytes.Buffer
	fc := newTestFileComparer(&out, 1)

	fc.recordMismatch(11, normalizedAction{Type: "reach", Actor: new(1)}, nil, "first mismatch")
	fc.recordMismatch(12, normalizedAction{Type: "hora", Actor: new(1), Target: new(0), Pai: "5m"}, nil, "second mismatch")

	if fc.fileSummary.mismatches != 2 {
		t.Errorf("mismatches = %d, want 2", fc.fileSummary.mismatches)
	}
	if got := strings.Count(out.String(), "mismatch:"); got != 1 {
		t.Errorf("reported mismatch details = %d, want 1; output = %q", got, out.String())
	}
	if strings.Contains(out.String(), "second mismatch") {
		t.Errorf("output = %q, want second mismatch suppressed", out.String())
	}
}

func TestSummaryAddIncludesImplicitPassesAndErrors(t *testing.T) {
	s := summary{files: 2, decisions: 3, matches: 1, implicitPasses: 1, mismatches: 1, errors: 1}
	s.add(summary{files: 4, decisions: 5, matches: 2, implicitPasses: 2, mismatches: 1, errors: 3})

	if s.files != 2 || s.decisions != 8 || s.matches != 3 || s.implicitPasses != 3 || s.mismatches != 2 || s.errors != 4 {
		t.Errorf("summary = %+v, want files unchanged and counters added", s)
	}
}

func newTestFileComparer(out *bytes.Buffer, limit int) *fileComparer {
	return &fileComparer{
		parent: &comparer{
			cfg: config{limit: limit},
			out: out,
		},
		path: "test.mjson",
		self: 1,
	}
}
