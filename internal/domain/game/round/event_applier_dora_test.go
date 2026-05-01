package round

import (
	"testing"

	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/event"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/tile"
)

func TestState_Apply_Dora_ReturnsErrorBeforeKan(t *testing.T) {
	s := mustNewRoundStateForTest(t, newValidHands())
	doraIndicator := *tile.MustTileFromCode("6p")

	if err := s.Apply(event.NewDora(doraIndicator)); err == nil {
		t.Fatal("Apply(Dora) succeeded unexpectedly")
	}

	if got := s.DoraIndicators(); len(got) != 1 {
		t.Fatalf("DoraIndicators() = %v, want unchanged initial indicator only", got)
	}
}
