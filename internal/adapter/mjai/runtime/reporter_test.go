package mjairuntime

import (
	"strings"
	"testing"
)

func TestReporter_ReportDecisionTrace(t *testing.T) {
	var out strings.Builder
	reporter := newReporter(&out)

	if err := reporter.ReportDecisionTrace("evaluation trace\n"); err != nil {
		t.Fatalf("ReportDecisionTrace() failed: %v", err)
	}

	if out.String() != "evaluation trace\n" {
		t.Errorf("output = %q, want evaluation trace", out.String())
	}
}

func TestReporter_ReportDecisionTrace_IgnoresEmptyTrace(t *testing.T) {
	var out strings.Builder
	reporter := newReporter(&out)

	if err := reporter.ReportDecisionTrace(""); err != nil {
		t.Fatalf("ReportDecisionTrace() failed: %v", err)
	}

	if out.String() != "" {
		t.Errorf("output = %q, want empty", out.String())
	}
}
