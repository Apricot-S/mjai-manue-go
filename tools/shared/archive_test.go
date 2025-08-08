package shared

import (
	"compress/gzip"
	"errors"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/Apricot-S/mjai-manue-go/internal/base"
	"github.com/Apricot-S/mjai-manue-go/internal/game"
	"github.com/Apricot-S/mjai-manue-go/internal/game/event/inbound"
	"github.com/Apricot-S/mjai-manue-go/internal/protocol"
	"github.com/Apricot-S/mjai-manue-go/internal/protocol/mjai"
)

func TestNewArchive(t *testing.T) {
	type args struct {
		paths   []string
		adapter protocol.Adapter
	}
	tests := []struct {
		name string
		args args
		want *Archive
	}{
		{
			name: "multiple paths",
			args: args{
				paths:   []string{"test1.json", "test1.json.gz"},
				adapter: &mjai.MjaiAdapter{},
			},
			want: &Archive{
				paths:   []string{"test1.json", "test1.json.gz"},
				adapter: &mjai.MjaiAdapter{},
				state:   &game.StateImpl{},
			},
		},
		{
			name: "empty",
			args: args{
				paths:   []string{},
				adapter: &mjai.MjaiAdapter{},
			},
			want: &Archive{
				paths:   []string{},
				adapter: &mjai.MjaiAdapter{},
				state:   &game.StateImpl{},
			},
		},
		{
			name: "nil",
			args: args{
				paths:   nil,
				adapter: &mjai.MjaiAdapter{},
			},
			want: &Archive{
				paths:   nil,
				adapter: &mjai.MjaiAdapter{},
				state:   &game.StateImpl{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewArchive(tt.args.paths, tt.args.adapter); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewArchive() = %v, want %v", got, tt.want)
			}
		})
	}
}

func writeTestFile(t *testing.T, name string, lines []string) string {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, name)
	content := strings.Join(lines, "\n")
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatalf("failed to write file: %v", err)
	}
	return path
}

func writeTestGZFile(t *testing.T, name string, plain []byte) string {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, name)

	file, err := os.Create(path)
	if err != nil {
		t.Fatal(err)
	}
	gz := gzip.NewWriter(file)
	if _, err := gz.Write(plain); err != nil {
		t.Fatal(err)
	}
	gz.Close()
	file.Close()
	return path
}

func TestArchive_PlayLight_SingleFile(t *testing.T) {
	data := []string{
		`{"type":"tsumo","actor":0,"pai":"F"}`,
		`{"type":"dahai","actor":0,"pai":"F","tsumogiri":true}`,
	}
	path := writeTestFile(t, "actions.json", data)

	pai, _ := base.NewPaiWithName("F")
	tsumo, _ := inbound.NewTsumo(0, *pai)
	dahai, _ := inbound.NewDahai(0, *pai, true)
	want := []inbound.Event{tsumo, dahai}

	archive := NewArchive([]string{path}, &mjai.MjaiAdapter{})
	var got []inbound.Event
	err := archive.PlayLight(func(act inbound.Event) error {
		got = append(got, act)
		return nil
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("unexpected actions: got %v, want %v", got, want)
	}
}

func TestArchive_PlayLight_GzipFile(t *testing.T) {
	plain := []byte(`{"type":"tsumo","actor":0,"pai":"F"}` + "\n" + `{"type":"dahai","actor":0,"pai":"F","tsumogiri":true}`)
	path := writeTestGZFile(t, "data.json.gz", plain)

	pai, _ := base.NewPaiWithName("F")
	tsumo, _ := inbound.NewTsumo(0, *pai)
	dahai, _ := inbound.NewDahai(0, *pai, true)
	want := []inbound.Event{tsumo, dahai}

	archive := NewArchive([]string{path}, &mjai.MjaiAdapter{})
	var got []inbound.Event
	err := archive.PlayLight(func(act inbound.Event) error {
		got = append(got, act)
		return nil
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("unexpected actions: got %v, want %v", got, want)
	}
}

func TestArchive_PlayLight_InvalidJSON(t *testing.T) {
	data := []string{
		`{"type":"tsumo","actor":0,"pai":"F"}`,
		`{"type":"dahai","actor":0`, // malformed JSON
	}
	path := writeTestFile(t, "broken.json", data)

	archive := NewArchive([]string{path}, &mjai.MjaiAdapter{})
	err := archive.PlayLight(func(act inbound.Event) error { return nil })

	if err == nil {
		t.Errorf("expected JSON error, got: %v", err)
	}
}

func TestArchive_PlayLight_FileNotFound(t *testing.T) {
	archive := NewArchive([]string{"/nonexistent/path.json"}, &mjai.MjaiAdapter{})
	err := archive.PlayLight(func(act inbound.Event) error { return nil })

	if err == nil || !strings.Contains(err.Error(), "failed to open") {
		t.Errorf("expected file open error, got: %v", err)
	}
}

func TestArchive_PlayLight_ErrorInCallback(t *testing.T) {
	data := []string{
		`{"type":"tsumo","actor":0,"pai":"F"}`,
		`{"type":"dahai","actor":0,"pai":"F","tsumogiri":true}`,
	}
	path := writeTestFile(t, "actions.json", data)

	archive := NewArchive([]string{path}, &mjai.MjaiAdapter{})
	err := archive.PlayLight(func(act inbound.Event) error {
		return errors.New("")
	})

	if err == nil || !strings.Contains(err.Error(), "failed to callback") {
		t.Errorf("expected callback error, got: %v", err)
	}
}
