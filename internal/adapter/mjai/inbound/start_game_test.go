package inbound_test

import (
	"reflect"
	"testing"

	"github.com/Apricot-S/mjai-manue-go/internal/adapter/mjai/inbound"
)

func TestParseMessage_StartGame(t *testing.T) {
	msg, err := inbound.ParseMessage([]byte(`{"type":"start_game","id":0,"names":["a","b","c","d"]}`))
	if err != nil {
		t.Fatalf("ParseMessage() failed: %v", err)
	}

	got, ok := msg.(*inbound.StartGame)
	if !ok {
		t.Fatalf("ParseMessage() = %T, want *inbound.StartGame", msg)
	}
	if got.Type != "start_game" {
		t.Errorf("Type = %q, want start_game", got.Type)
	}
	if got.ID == nil {
		t.Fatal("ID = nil, want 0")
	}
	if *got.ID != 0 {
		t.Errorf("ID = %d, want 0", *got.ID)
	}
	if want := []string{"a", "b", "c", "d"}; !reflect.DeepEqual(got.Names, want) {
		t.Errorf("Names = %v, want %v", got.Names, want)
	}
}

func TestParseMessage_StartGame_NoNames(t *testing.T) {
	msg, err := inbound.ParseMessage([]byte(`{"type":"start_game","id":0}`))
	if err != nil {
		t.Fatalf("ParseMessage() failed: %v", err)
	}

	got, ok := msg.(*inbound.StartGame)
	if !ok {
		t.Fatalf("ParseMessage() = %T, want *inbound.StartGame", msg)
	}
	if got.Names != nil {
		t.Errorf("Names = %v, want nil", got.Names)
	}
}

func TestParseMessage_StartGame_NoID(t *testing.T) {
	msg, err := inbound.ParseMessage([]byte(`{"type":"start_game","names":["a","b","c","d"]}`))
	if err != nil {
		t.Fatalf("ParseMessage() failed: %v", err)
	}

	got, ok := msg.(*inbound.StartGame)
	if !ok {
		t.Fatalf("ParseMessage() = %T, want *inbound.StartGame", msg)
	}
	if got.ID != nil {
		t.Errorf("ID = %d, want nil", *got.ID)
	}
	if want := []string{"a", "b", "c", "d"}; !reflect.DeepEqual(got.Names, want) {
		t.Errorf("Names = %v, want %v", got.Names, want)
	}
}
