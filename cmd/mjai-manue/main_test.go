package main

import (
	"strings"
	"testing"
)

func TestRun_HelloWritesJoin(t *testing.T) {
	in := strings.NewReader(`{"type":"hello","protocol":"mjsonp","protocol_version":3}` + "\n")
	var out strings.Builder
	var errOut strings.Builder

	got := run(nil, in, &out, &errOut)
	if got != exitOK {
		t.Fatalf("run() = %d, want %d; stderr = %q", got, exitOK, errOut.String())
	}

	want := `{"type":"join","name":"Manue030","room":"default"}` + "\n"
	if out.String() != want {
		t.Errorf("stdout = %q, want %q", out.String(), want)
	}
}

func TestRun_NameFlag(t *testing.T) {
	in := strings.NewReader(`{"type":"hello","protocol":"mjsonp","protocol_version":3}` + "\n")
	var out strings.Builder
	var errOut strings.Builder

	got := run([]string{"--name", "custom", "--seed", "123"}, in, &out, &errOut)
	if got != exitOK {
		t.Fatalf("run() = %d, want %d; stderr = %q", got, exitOK, errOut.String())
	}

	want := `{"type":"join","name":"custom","room":"default"}` + "\n"
	if out.String() != want {
		t.Errorf("stdout = %q, want %q", out.String(), want)
	}
}

func TestRun_IDFlagAllowsStartGameWithoutID(t *testing.T) {
	in := strings.NewReader(`{"type":"start_game"}` + "\n")
	var out strings.Builder
	var errOut strings.Builder

	got := run([]string{"--id", "2"}, in, &out, &errOut)
	if got != exitOK {
		t.Fatalf("run() = %d, want %d; stderr = %q", got, exitOK, errOut.String())
	}
	if out.String() != "" {
		t.Errorf("stdout = %q, want empty", out.String())
	}
}

func TestRun_InvalidIDFlagReturnsUsageError(t *testing.T) {
	var out strings.Builder
	var errOut strings.Builder

	got := run([]string{"--id", "4"}, strings.NewReader(""), &out, &errOut)
	if got != exitUsageError {
		t.Fatalf("run() = %d, want %d; stderr = %q", got, exitUsageError, errOut.String())
	}
	if out.String() != "" {
		t.Errorf("stdout = %q, want empty", out.String())
	}
	if !strings.Contains(errOut.String(), "invalid player seat: 4") {
		t.Errorf("stderr = %q, want invalid player seat", errOut.String())
	}
}

func TestRun_TooManyArguments(t *testing.T) {
	var out strings.Builder
	var errOut strings.Builder

	got := run([]string{"mjsonp://127.0.0.1:11600/default", "extra"}, strings.NewReader(""), &out, &errOut)
	if got != exitUsageError {
		t.Errorf("run() = %d, want %d", got, exitUsageError)
	}
	if out.String() != "" {
		t.Errorf("stdout = %q, want empty", out.String())
	}
}
