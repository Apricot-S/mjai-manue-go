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

func TestRun_IDFlagAllowsStartGameWithoutID(t *testing.T) {
	var stdout strings.Builder
	var stderr strings.Builder

	code := run([]string{"--id", "2"}, strings.NewReader(`{"type":"start_game"}`+"\n"), &stdout, &stderr)
	if code != exitOK {
		t.Fatalf("run() = %d, want %d; stderr = %q", code, exitOK, stderr.String())
	}
	if stdout.String() != "" {
		t.Errorf("stdout = %q, want empty", stdout.String())
	}
}

func TestRun_InvalidIDFlagReturnsUsageError(t *testing.T) {
	var stdout strings.Builder
	var stderr strings.Builder

	code := run([]string{"--id", "4"}, strings.NewReader(""), &stdout, &stderr)
	if code != exitUsageError {
		t.Fatalf("run() = %d, want %d; stderr = %q", code, exitUsageError, stderr.String())
	}
	if stdout.String() != "" {
		t.Errorf("stdout = %q, want empty", stdout.String())
	}
	if !strings.Contains(stderr.String(), "invalid player seat: 4") {
		t.Errorf("stderr = %q, want invalid player seat", stderr.String())
	}
}
