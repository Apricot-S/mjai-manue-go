package round

import (
	"testing"

	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/event"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/tile"
)

func TestState_Apply_Dora(t *testing.T) {
	s := mustNewRoundStateForTest(t, newValidHands())
	doraIndicator := *tile.MustTileFromCode("6p")

	if err := s.Apply(event.NewDora(doraIndicator)); err != nil {
		t.Fatalf("Apply(Dora) failed: %v", err)
	}

	if got := s.DoraIndicators(); len(got) != 2 || got[1] != doraIndicator {
		t.Fatalf("DoraIndicators() = %v, want appended %v", got, doraIndicator)
	}
}
