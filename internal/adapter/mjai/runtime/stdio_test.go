package runtime_test

import (
	"strings"
	"testing"

	mjairuntime "github.com/Apricot-S/mjai-manue-go/internal/adapter/mjai/runtime"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/ai"
)

func TestRunStdio_HelloWritesJoin(t *testing.T) {
	in := strings.NewReader(`{"type":"hello","protocol":"mjsonp","protocol_version":3}` + "\n")
	var out strings.Builder

	err := mjairuntime.RunStdio(mjairuntime.StdioConfig{
		Name:  "tsumogiri",
		Room:  "default",
		Agent: ai.NewTsumogiriAgent(),
		In:    in,
		Out:   &out,
	})
	if err != nil {
		t.Fatalf("RunStdio() failed: %v", err)
	}

	want := `{"type":"join","name":"tsumogiri","room":"default"}` + "\n"
	if out.String() != want {
		t.Errorf("output = %q, want %q", out.String(), want)
	}
}

func TestRunStdio_EmptyLine(t *testing.T) {
	in := strings.NewReader("\n")
	var out strings.Builder

	err := mjairuntime.RunStdio(mjairuntime.StdioConfig{
		Name:  "tsumogiri",
		Room:  "default",
		Agent: ai.NewTsumogiriAgent(),
		In:    in,
		Out:   &out,
	})
	if err == nil {
		t.Fatal("RunStdio() succeeded unexpectedly")
	}
}
