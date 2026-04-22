package inbound_test

import (
	"testing"

	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/event"
	"github.com/Apricot-S/mjai-manue-go/internal/infrastructure/mjai/inbound"
)

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

		got, err := inbound.ParseEvent([]byte(payload))
		if err != nil {
			t.Fatalf("ParseEvent() failed: %v", err)
		}
		if _, ok := got.(*event.StartRound); !ok {
			t.Fatalf("ParseEvent() = %T, want *event.StartRound", got)
		}
	})

	t.Run("tsumo", func(t *testing.T) {
		got, err := inbound.ParseEvent([]byte(`{"type":"tsumo","actor":1,"pai":"E"}`))
		if err != nil {
			t.Fatalf("ParseEvent() failed: %v", err)
		}
		if _, ok := got.(*event.Draw); !ok {
			t.Fatalf("ParseEvent() = %T, want *event.Draw", got)
		}
	})

	t.Run("unknown type", func(t *testing.T) {
		if _, err := inbound.ParseEvent([]byte(`{"type":"nope"}`)); err == nil {
			t.Fatal("ParseEvent() succeeded unexpectedly")
		}
	})

	t.Run("missing type", func(t *testing.T) {
		if _, err := inbound.ParseEvent([]byte(`{"actor":1,"pai":"E"}`)); err == nil {
			t.Fatal("ParseEvent() succeeded unexpectedly")
		}
	})
}
