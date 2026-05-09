package mjairuntime

import (
	"bytes"
	"encoding/json/v2"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/Apricot-S/mjai-manue-go/internal/domain/ai"
)

func TestGoldenStdout(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		golden string
		policy jsonLinesPolicy
	}{
		{
			name:   "stdio",
			input:  "testdata/tsumogiri/self_draw.input.mjson",
			golden: "testdata/tsumogiri/self_draw.stdio.golden",
			policy: jsonLinesPolicy{},
		},
		{
			name:   "mjsonp",
			input:  "testdata/tsumogiri/self_draw.input.mjson",
			golden: "testdata/tsumogiri/self_draw.mjsonp.golden",
			policy: jsonLinesPolicy{
				respondNoneOnNoReaction: true,
				stopOnEndGame:           true,
				errorOnEOFBeforeEndGame: true,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			input := readGoldenFile(t, tt.input)
			want := readGoldenFile(t, tt.golden)

			var out bytes.Buffer
			err := runJSONLines("tsumogiri", "default", ai.NewTsumogiriAgent(), strings.NewReader(input), &out, nil, tt.policy)
			if err != nil {
				t.Fatalf("runJSONLines() failed: %v", err)
			}

			gotStdout := out.String()
			if strings.HasSuffix(gotStdout, "\n") != strings.HasSuffix(want, "\n") {
				t.Fatalf("stdout trailing newline mismatch: got %t, want %t", strings.HasSuffix(gotStdout, "\n"), strings.HasSuffix(want, "\n"))
			}
			gotMessages, err := jsonMessages(gotStdout)
			if err != nil {
				t.Fatalf("jsonMessages(stdout) failed: %v", err)
			}
			wantMessages, err := jsonMessages(want)
			if err != nil {
				t.Fatalf("jsonMessages(golden) failed: %v", err)
			}
			if !reflect.DeepEqual(gotMessages, wantMessages) {
				t.Errorf("stdout =\n%s\nwant\n%s", formatJSONMessages(t, gotMessages), formatJSONMessages(t, wantMessages))
			}
		})
	}
}

func readGoldenFile(t *testing.T, name string) string {
	t.Helper()

	b, err := os.ReadFile(filepath.FromSlash(name))
	if err != nil {
		t.Fatalf("ReadFile(%q) failed: %v", name, err)
	}
	return string(b)
}

func jsonMessages(s string) ([]map[string]any, error) {
	var messages []map[string]any

	// JSON Lines output normally ends with one trailing newline. Remove only
	// that final line terminator; blank lines elsewhere must remain visible.
	lines := strings.Split(s, "\n")
	if len(lines) > 0 && lines[len(lines)-1] == "" {
		lines = lines[:len(lines)-1]
	}
	for _, line := range lines {
		// Accept CRLF fixtures and runtime output without treating '\r' as JSON.
		if before, ok := strings.CutSuffix(line, "\r"); ok {
			line = before
		}

		// Decode each physical line. Empty or whitespace-only lines are not
		// skipped, so unexpected stdout is reported as a JSON parse failure.
		var msg map[string]any
		if err := json.Unmarshal([]byte(line), &msg); err != nil {
			return nil, err
		}

		// Store decoded objects so JSON field order differences are allowed.
		messages = append(messages, msg)
	}
	return messages, nil
}

func formatJSONMessages(t *testing.T, messages []map[string]any) string {
	t.Helper()

	var b strings.Builder
	for _, msg := range messages {
		line, err := json.Marshal(msg)
		if err != nil {
			t.Fatalf("Marshal(%v) failed: %v", msg, err)
		}
		b.Write(line)
		b.WriteByte('\n')
	}
	return b.String()
}
