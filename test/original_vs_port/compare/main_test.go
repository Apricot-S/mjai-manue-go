package main

import (
	"bytes"
	"strings"
	"testing"
)

func TestRun_UsageErrorExitCode(t *testing.T) {
	var out bytes.Buffer
	var errOut bytes.Buffer
	got := run(nil, &out, &errOut)
	if got != exitRunError {
		t.Errorf("exit code = %d, want %d", got, exitRunError)
	}
	if !strings.Contains(errOut.String(), "usage:") {
		t.Errorf("stderr = %q, want usage", errOut.String())
	}
}
