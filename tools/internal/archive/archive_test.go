package archive

import (
	"compress/gzip"
	"os"
	"path/filepath"
	"testing"

	"github.com/Apricot-S/mjai-manue-go/internal/adapter/mjai/inbound"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/event"
)

const sampleLog = `{"type":"start_game","names":["a","b","c","d"]}
{"type":"start_kyoku","bakaze":"E","kyoku":1,"honba":0,"kyotaku":0,"oya":0,"dora_marker":"1m","tehais":[["2m","3m","4m","5m","6m","7m","8m","9m","1p","2p","3p","4p","5p"],["1m","1m","2m","2m","3m","3m","4m","4m","5m","5m","6m","6m","7m"],["1s","2s","3s","4s","5s","6s","7s","8s","9s","E","S","W","N"],["1p","1p","2p","2p","3p","3p","4p","4p","5p","5p","6p","6p","7p"]],"scores":[25000,25000,25000,25000]}
{"type":"tsumo","actor":0,"pai":"1m"}
{"type":"dahai","actor":0,"pai":"1m","tsumogiri":true}
{"type":"end_kyoku"}
{"type":"end_game"}
`

func TestArchivePlayParsesMessagesAndUpdatesState(t *testing.T) {
	path := writeTempFile(t, "sample.mjson", sampleLog)
	archive := NewArchive()

	var numMessages int
	var donePaths []string
	var eventTypes []string
	var turnsAfterDiscard []float64
	err := archive.PlayPaths([]string{path}, Handlers{
		OnMessage: func(msg inbound.Message) error {
			numMessages++
			return nil
		},
		OnEvent: func(ev event.Event, archive *Archive) error {
			eventTypes = append(eventTypes, eventTypeName(ev))
			if _, ok := ev.(*event.Discard); ok {
				state, ok := archive.StateViewer()
				if !ok {
					t.Error("StateViewer() missing after discard")
					return nil
				}
				turnsAfterDiscard = append(turnsAfterDiscard, state.Turn())
			}
			return nil
		},
		OnFileDone: func(path string) error {
			donePaths = append(donePaths, path)
			return nil
		},
	})
	if err != nil {
		t.Fatalf("Archive.PlayPaths() error = %v", err)
	}

	if numMessages != 6 {
		t.Errorf("num messages = %d, want 6", numMessages)
	}
	wantEventTypes := []string{"start_round", "draw", "discard", "end_round"}
	if len(eventTypes) != len(wantEventTypes) {
		t.Fatalf("event types = %v, want %v", eventTypes, wantEventTypes)
	}
	for i := range eventTypes {
		if eventTypes[i] != wantEventTypes[i] {
			t.Errorf("eventTypes[%d] = %q, want %q", i, eventTypes[i], wantEventTypes[i])
		}
	}
	if len(turnsAfterDiscard) != 1 || turnsAfterDiscard[0] != 0.25 {
		t.Errorf("turns after discard = %v, want [0.25]", turnsAfterDiscard)
	}
	if len(donePaths) != 1 || donePaths[0] != path {
		t.Errorf("done paths = %v, want [%s]", donePaths, path)
	}
	if _, ok := archive.State(); ok {
		t.Error("State() exists after end_kyoku/end_game")
	}
}

func TestArchivePlayReadsGzip(t *testing.T) {
	path := writeTempGzipFile(t, "sample.mjson.gz", sampleLog)
	archive := NewArchive()

	var numEvents int
	if err := archive.PlayPaths([]string{path}, Handlers{
		OnEvent: func(event.Event, *Archive) error {
			numEvents++
			return nil
		},
	}); err != nil {
		t.Fatalf("Archive.PlayPaths() error = %v", err)
	}
	if numEvents != 4 {
		t.Errorf("num events = %d, want 4", numEvents)
	}
}

func TestArchivePlayRejectsEmptyLine(t *testing.T) {
	path := writeTempFile(t, "empty.mjson", "{}\n\n")
	archive := NewArchive()

	if err := archive.PlayPaths([]string{path}, Handlers{}); err == nil {
		t.Fatal("Archive.PlayPaths() succeeded unexpectedly")
	}
}

func TestGlobAll(t *testing.T) {
	dir := t.TempDir()
	writeTempFileAt(t, filepath.Join(dir, "a.mjson"), "{}\n")
	writeTempFileAt(t, filepath.Join(dir, "b.txt"), "{}\n")

	got, err := GlobAll([]string{filepath.Join(dir, "*.mjson")})
	if err != nil {
		t.Fatalf("GlobAll() error = %v", err)
	}
	if len(got) != 1 || filepath.Base(got[0]) != "a.mjson" {
		t.Errorf("GlobAll() = %v, want only a.mjson", got)
	}
}

func eventTypeName(ev event.Event) string {
	switch ev.(type) {
	case *event.StartRound:
		return "start_round"
	case *event.Draw:
		return "draw"
	case *event.Discard:
		return "discard"
	case *event.EndRound:
		return "end_round"
	default:
		return "unknown"
	}
}

func writeTempFile(t *testing.T, name string, content string) string {
	t.Helper()
	return writeTempFileAt(t, filepath.Join(t.TempDir(), name), content)
}

func writeTempFileAt(t *testing.T, path string, content string) string {
	t.Helper()
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatalf("failed to write temp file: %v", err)
	}
	return path
}

func writeTempGzipFile(t *testing.T, name string, content string) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), name)
	file, err := os.Create(path)
	if err != nil {
		t.Fatalf("failed to create gzip file: %v", err)
	}
	gz := gzip.NewWriter(file)
	if _, err := gz.Write([]byte(content)); err != nil {
		t.Fatalf("failed to write gzip content: %v", err)
	}
	if err := gz.Close(); err != nil {
		t.Fatalf("failed to close gzip writer: %v", err)
	}
	if err := file.Close(); err != nil {
		t.Fatalf("failed to close gzip file: %v", err)
	}
	return path
}
