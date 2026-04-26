package outbound_test

import (
	"testing"

	"github.com/Apricot-S/mjai-manue-go/internal/adapter/mjai/outbound"
)

func TestMarshalMessage_Join(t *testing.T) {
	got, err := outbound.MarshalMessage(outbound.NewJoin("A", "default"))
	if err != nil {
		t.Fatalf("MarshalMessage() failed: %v", err)
	}
	if want := `{"type":"join","name":"A","room":"default"}`; string(got) != want {
		t.Errorf("MarshalMessage() = %s, want %s", got, want)
	}
}

func TestMarshalMessage_Join_EmptyFields(t *testing.T) {
	got, err := outbound.MarshalMessage(outbound.NewJoin("", ""))
	if err != nil {
		t.Fatalf("MarshalMessage() failed: %v", err)
	}
	if want := `{"type":"join","name":"","room":""}`; string(got) != want {
		t.Errorf("MarshalMessage() = %s, want %s", got, want)
	}
}
