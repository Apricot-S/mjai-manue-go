package main

import (
	"strings"
	"testing"
)

func TestRun_InvalidURLReturnsUsageError(t *testing.T) {
	var stdout strings.Builder
	var stderr strings.Builder

	code := run([]string{"stdio://127.0.0.1:11600/room"}, strings.NewReader(""), &stdout, &stderr)
	if code != exitUsageError {
		t.Errorf("run() = %d, want %d", code, exitUsageError)
	}
	if stdout.String() != "" {
		t.Errorf("stdout = %q, want empty", stdout.String())
	}
	if !strings.Contains(stderr.String(), "unsupported URL scheme") {
		t.Errorf("stderr = %q, want unsupported URL scheme", stderr.String())
	}
}

func TestRun_TooManyArgumentsReturnsUsageError(t *testing.T) {
	var stdout strings.Builder
	var stderr strings.Builder

	code := run([]string{"mjsonp://127.0.0.1:11600/room", "extra"}, strings.NewReader(""), &stdout, &stderr)
	if code != exitUsageError {
		t.Errorf("run() = %d, want %d", code, exitUsageError)
	}
}
