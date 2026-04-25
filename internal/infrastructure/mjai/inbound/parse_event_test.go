package inbound_test

import (
	"testing"

	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/event"
	"github.com/Apricot-S/mjai-manue-go/internal/infrastructure/mjai/inbound"
)

func TestParseEvent_FromMessage(t *testing.T) {
	got, err := inbound.ParseEvent(&inbound.Tsumo{
		Type:  "tsumo",
		Actor: 1,
		Pai:   "E",
	})
	if err != nil {
		t.Fatalf("ParseEvent() failed: %v", err)
	}
	if _, ok := got.(*event.Draw); !ok {
		t.Fatalf("ParseEvent() = %T, want *event.Draw", got)
	}
}

func TestParseEvent_NonEventMessage(t *testing.T) {
	if _, err := inbound.ParseEvent(&inbound.Hello{Type: "hello"}); err == nil {
		t.Fatal("ParseEvent() succeeded unexpectedly")
	}
}

func TestParseEvent_Dispatch(t *testing.T) {
	t.Run("start_kyoku", func(t *testing.T) {
		payload := `{
			"type":"start_kyoku",
			"bakaze":"E",
			"kyoku":1,
			"honba":0,
			"kyotaku":0,
			"oya":0,
			"dora_marker":"5mr",
			"tehais":[` +
			toJSONHand(unknownHand()) + "," +
			toJSONHand(unknownHand()) + "," +
			toJSONHand(unknownHand()) + "," +
			toJSONHand(unknownHand()) + `]
		}`

		msg, err := inbound.ParseMessage([]byte(payload))
		if err != nil {
			t.Fatalf("ParseMessage() failed: %v", err)
		}
		got, err := inbound.ParseEvent(msg)
		if err != nil {
			t.Fatalf("ParseEvent() failed: %v", err)
		}
		if _, ok := got.(*event.StartRound); !ok {
			t.Fatalf("ParseEvent() = %T, want *event.StartRound", got)
		}
	})

	t.Run("tsumo", func(t *testing.T) {
		msg, err := inbound.ParseMessage([]byte(`{"type":"tsumo","actor":1,"pai":"E"}`))
		if err != nil {
			t.Fatalf("ParseMessage() failed: %v", err)
		}
		got, err := inbound.ParseEvent(msg)
		if err != nil {
			t.Fatalf("ParseEvent() failed: %v", err)
		}
		if _, ok := got.(*event.Draw); !ok {
			t.Fatalf("ParseEvent() = %T, want *event.Draw", got)
		}
	})

	t.Run("dahai", func(t *testing.T) {
		msg, err := inbound.ParseMessage([]byte(`{"type":"dahai","actor":1,"pai":"W","tsumogiri":false}`))
		if err != nil {
			t.Fatalf("ParseMessage() failed: %v", err)
		}
		got, err := inbound.ParseEvent(msg)
		if err != nil {
			t.Fatalf("ParseEvent() failed: %v", err)
		}
		if _, ok := got.(*event.Discard); !ok {
			t.Fatalf("ParseEvent() = %T, want *event.Discard", got)
		}
	})

	t.Run("end_kyoku", func(t *testing.T) {
		msg, err := inbound.ParseMessage([]byte(`{"type":"end_kyoku"}`))
		if err != nil {
			t.Fatalf("ParseMessage() failed: %v", err)
		}
		got, err := inbound.ParseEvent(msg)
		if err != nil {
			t.Fatalf("ParseEvent() failed: %v", err)
		}
		if _, ok := got.(*event.EndRound); !ok {
			t.Fatalf("ParseEvent() = %T, want *event.EndRound", got)
		}
	})
}
