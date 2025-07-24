package shared

import (
	"compress/gzip"
	"errors"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/Apricot-S/mjai-manue-go/internal/game"
	"github.com/go-json-experiment/json/jsontext"
)

func TestNewArchive(t *testing.T) {
	type args struct {
		paths []string
	}
	tests := []struct {
		name string
		args args
		want *Archive
	}{
		{
			name: "multiple paths",
			args: args{paths: []string{"test1.json", "test1.json.gz"}},
			want: &Archive{
				paths: []string{"test1.json", "test1.json.gz"},
				state: &game.StateImpl{},
			},
		},
		{
			name: "empty",
			args: args{paths: []string{}},
			want: &Archive{
				paths: []string{},
				state: &game.StateImpl{},
			},
		},
		{
			name: "nil",
			args: args{paths: nil},
			want: &Archive{
				paths: nil,
				state: &game.StateImpl{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewArchive(tt.args.paths); !reflect.DeepEqual(got, tt.want) {
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
	data := []string{`{"key":"value1"}`, `{"key":"value2"}`}
	path := writeTestFile(t, "actions.json", data)

	want := []jsontext.Value{}
	for _, d := range data {
		want = append(want, jsontext.Value(d))
	}

	archive := NewArchive([]string{path})
	var got []jsontext.Value
	err := archive.PlayLight(func(act jsontext.Value) error {
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
	plain := []byte(`{"x":1}` + "\n" + `{"x":2}`)
	path := writeTestGZFile(t, "data.json.gz", plain)

	want := []jsontext.Value{
		jsontext.Value(`{"x":1}`),
		jsontext.Value(`{"x":2}`),
	}

	archive := NewArchive([]string{path})
	var got []jsontext.Value
	err := archive.PlayLight(func(act jsontext.Value) error {
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
		`{"valid":true}`,
		`{invalid`, // malformed JSON
	}
	path := writeTestFile(t, "broken.json", data)

	archive := NewArchive([]string{path})
	err := archive.PlayLight(func(act jsontext.Value) error { return nil })

	if err == nil || !strings.Contains(err.Error(), "json decode error") {
		t.Errorf("expected JSON error, got: %v", err)
	}
}

func TestArchive_PlayLight_FileNotFound(t *testing.T) {
	archive := NewArchive([]string{"/nonexistent/path.json"})
	err := archive.PlayLight(func(act jsontext.Value) error { return nil })

	if err == nil || !strings.Contains(err.Error(), "failed to open") {
		t.Errorf("expected file open error, got: %v", err)
	}
}

func TestArchive_PlayLight_ErrorInCallback(t *testing.T) {
	data := []string{`{"key":"value1"}`, `{"key":"value2"}`}
	path := writeTestFile(t, "actions.json", data)

	archive := NewArchive([]string{path})
	err := archive.PlayLight(func(act jsontext.Value) error {
		return errors.New("")
	})

	if err == nil || !strings.Contains(err.Error(), "failed to callback") {
		t.Errorf("expected callback error, got: %v", err)
	}
}
