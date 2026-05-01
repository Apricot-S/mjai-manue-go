package round

import (
	"testing"

	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/common"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/event"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/seat"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/tile"
)

func TestState_Apply_Draw(t *testing.T) {
	s := mustNewRoundStateForTest(t, newValidHands())
	actor := *seat.MustSeat(0)
	drawnTile := *tile.MustTileFromCode("6m")

	before := s.NumLeftTiles()
	if err := s.Apply(event.NewDraw(actor, drawnTile)); err != nil {
		t.Fatalf("Apply(Draw) failed: %v", err)
	}

	if got := s.NumLeftTiles(); got != before-1 {
		t.Errorf("NumLeftTiles() = %d, want %d", got, before-1)
	}
	if got := s.Player(actor).DrawnTile(); got == nil || got.ID() != drawnTile.ID() {
		t.Fatalf("DrawnTile() = %v, want %v", got, drawnTile)
	}
	if !s.Player(actor).CanDiscard() {
		t.Error("CanDiscard() = false, want true")
	}

	for i := 1; i < common.NumPlayers; i++ {
		playerSeat := *seat.MustSeat(i)
		if got := s.Player(playerSeat).DrawnTile(); got != nil {
			t.Errorf("player %d DrawnTile() = %v, want nil", i, got)
		}
	}
}
