package round

import (
	"testing"

	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/common"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/event"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/round/player"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/seat"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/tile"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/wind"
)

func TestState_Apply_Draw(t *testing.T) {
	s := mustNewRoundStateForTest(t, newValidHands())
	actor := seat.MustSeat(0)
	drawnTile := tile.MustTileFromCode("6m")

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
		playerSeat := seat.MustSeat(i)
		if got := s.Player(playerSeat).DrawnTile(); got != nil {
			t.Errorf("player %d DrawnTile() = %v, want nil", i, got)
		}
	}
}

func TestState_Apply_Draw_ReturnsErrorWhenActorIsNotDealerAtRoundStart(t *testing.T) {
	s := mustNewRoundStateForTest(t, newValidHands())
	actor := seat.MustSeat(1)
	drawnTile := tile.MustTileFromCode("6m")

	if err := s.Apply(event.NewDraw(actor, drawnTile)); err == nil {
		t.Fatal("Apply(Draw) succeeded unexpectedly")
	}

	if got := s.NumLeftTiles(); got != NumInitWall {
		t.Errorf("NumLeftTiles() = %d, want %d", got, NumInitWall)
	}
	if got := s.Player(actor).DrawnTile(); got != nil {
		t.Fatalf("DrawnTile() = %v, want nil", got)
	}
}

func TestState_Apply_Draw_ReturnsErrorWhenNoTilesLeft(t *testing.T) {
	actorPlayer, err := player.NewVisiblePlayer(newValidHands()[0])
	if err != nil {
		t.Fatalf("player.NewVisiblePlayer() failed: %v", err)
	}
	s := NewStateForTest(
		wind.East,
		1,
		0,
		0,
		[common.NumPlayers]int{25000, 25000, 25000, 25000},
		seat.MustSeat(0),
		seat.MustSeat(0),
		tile.Tiles{tile.MustTileFromCode("E")},
		0,
		[common.NumPlayers]player.Player{
			actorPlayer,
			player.NewInvisiblePlayer(),
			player.NewInvisiblePlayer(),
			player.NewInvisiblePlayer(),
		},
	)
	actor := seat.MustSeat(0)
	drawnTile := tile.MustTileFromCode("6m")

	if err := s.Apply(event.NewDraw(actor, drawnTile)); err == nil {
		t.Fatal("Apply(Draw) succeeded unexpectedly")
	}

	if got := s.NumLeftTiles(); got != 0 {
		t.Errorf("NumLeftTiles() = %d, want 0", got)
	}
	if got := s.Player(actor).DrawnTile(); got != nil {
		t.Fatalf("DrawnTile() = %v, want nil", got)
	}
}
