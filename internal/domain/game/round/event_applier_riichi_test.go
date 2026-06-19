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

func riichiReadyHandForTest() [common.InitHandSize]tile.Tile {
	return [common.InitHandSize]tile.Tile{
		tile.MustTileFromCode("1m"), tile.MustTileFromCode("2m"), tile.MustTileFromCode("3m"),
		tile.MustTileFromCode("4p"), tile.MustTileFromCode("5p"), tile.MustTileFromCode("6p"),
		tile.MustTileFromCode("7s"), tile.MustTileFromCode("8s"), tile.MustTileFromCode("9s"),
		tile.MustTileFromCode("E"), tile.MustTileFromCode("E"), tile.MustTileFromCode("S"),
		tile.MustTileFromCode("W"),
	}
}

func TestState_Apply_Riichi(t *testing.T) {
	hands := newValidHands()
	hands[0] = riichiReadyHandForTest()
	s := mustNewRoundStateForTest(t, hands)
	actor := seat.MustSeat(0)

	if err := s.Apply(event.NewDraw(actor, tile.MustTileFromCode("S"))); err != nil {
		t.Fatalf("Apply(Draw) failed: %v", err)
	}
	if err := s.Apply(event.NewRiichi(actor)); err != nil {
		t.Fatalf("Apply(Riichi) failed: %v", err)
	}

	if got := s.Player(actor).RiichiState(); got != player.RiichiDeclared {
		t.Fatalf("RiichiState() = %v, want %v", got, player.RiichiDeclared)
	}
}

func TestState_Apply_Riichi_ReturnsErrorWithoutNextDrawTurn(t *testing.T) {
	hands := newValidHands()
	hands[0] = riichiReadyHandForTest()
	s := NewStateForTest(
		wind.East,
		1,
		0,
		0,
		[common.NumPlayers]int{25000, 25000, 25000, 25000},
		seat.MustSeat(0),
		seat.MustSeat(0),
		tile.Tiles{tile.MustTileFromCode("E")},
		common.NumPlayers,
		newVisiblePlayersForTest(t, hands),
	)
	actor := seat.MustSeat(0)

	if err := s.Apply(event.NewDraw(actor, tile.MustTileFromCode("S"))); err != nil {
		t.Fatalf("Apply(Draw) failed: %v", err)
	}
	if got := s.NumLeftTiles(); got != common.NumPlayers-1 {
		t.Fatalf("NumLeftTiles() = %d, want %d", got, common.NumPlayers-1)
	}

	if err := s.Apply(event.NewRiichi(actor)); err == nil {
		t.Fatal("Apply(Riichi) succeeded unexpectedly")
	}

	if got := s.Player(actor).RiichiState(); got != player.NotRiichi {
		t.Fatalf("RiichiState() = %v, want %v", got, player.NotRiichi)
	}
}

func TestState_Apply_Riichi_ReturnsErrorWhenActorIsNotPendingDiscardPlayer(t *testing.T) {
	hands := newValidHands()
	hands[0] = riichiReadyHandForTest()
	s := mustNewRoundStateForTest(t, hands)
	drawActor := seat.MustSeat(0)
	riichiActor := seat.MustSeat(1)

	if err := s.Apply(event.NewDraw(drawActor, tile.MustTileFromCode("S"))); err != nil {
		t.Fatalf("Apply(Draw) failed: %v", err)
	}
	if err := s.Apply(event.NewRiichi(riichiActor)); err == nil {
		t.Fatal("Apply(Riichi) succeeded unexpectedly")
	}

	if got := s.Player(riichiActor).RiichiState(); got != player.NotRiichi {
		t.Fatalf("RiichiState() = %v, want %v", got, player.NotRiichi)
	}
}

func TestState_Apply_DiscardAfterRiichi_ReturnsErrorWhenActorIsNotPendingDiscardPlayer(t *testing.T) {
	hands := newValidHands()
	hands[0] = riichiReadyHandForTest()
	s := mustNewRoundStateForTest(t, hands)
	drawActor := seat.MustSeat(0)
	discardActor := seat.MustSeat(1)
	drawnTile := tile.MustTileFromCode("S")
	discardedTile := tile.MustTileFromCode("W")

	if err := s.Apply(event.NewDraw(drawActor, drawnTile)); err != nil {
		t.Fatalf("Apply(Draw) failed: %v", err)
	}
	if err := s.Apply(event.NewRiichi(drawActor)); err != nil {
		t.Fatalf("Apply(Riichi) failed: %v", err)
	}
	if err := s.Apply(event.NewDiscard(discardActor, discardedTile, false)); err == nil {
		t.Fatal("Apply(Discard) succeeded unexpectedly")
	}

	if got := s.Player(drawActor).RiichiState(); got != player.RiichiDeclared {
		t.Fatalf("draw actor RiichiState() = %v, want %v", got, player.RiichiDeclared)
	}
	if got := s.Player(discardActor).River(); len(got) != 0 {
		t.Fatalf("discard actor River() = %v, want empty", got)
	}
}

func TestState_Apply_RiichiAccepted(t *testing.T) {
	tests := []struct {
		name       string
		deltas     *[common.NumPlayers]int
		scores     *[common.NumPlayers]int
		wantScores [common.NumPlayers]int
	}{
		{
			name:       "scores",
			scores:     &[common.NumPlayers]int{24000, 25000, 25000, 25000},
			wantScores: [common.NumPlayers]int{24000, 25000, 25000, 25000},
		},
		{
			name:       "deltas",
			deltas:     &[common.NumPlayers]int{-1000, 0, 0, 0},
			wantScores: [common.NumPlayers]int{24000, 25000, 25000, 25000},
		},
		{
			name:       "scores take precedence over deltas",
			deltas:     &[common.NumPlayers]int{1, 2, 3, 4},
			scores:     &[common.NumPlayers]int{24000, 25000, 25000, 25000},
			wantScores: [common.NumPlayers]int{24000, 25000, 25000, 25000},
		},
		{
			name:       "no scores or deltas",
			wantScores: [common.NumPlayers]int{24000, 25000, 25000, 25000},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hands := newValidHands()
			hands[0] = riichiReadyHandForTest()
			s := mustNewRoundStateForTest(t, hands)
			actor := seat.MustSeat(0)

			if err := s.Apply(event.NewDraw(actor, tile.MustTileFromCode("S"))); err != nil {
				t.Fatalf("Apply(Draw) failed: %v", err)
			}
			if err := s.Apply(event.NewRiichi(actor)); err != nil {
				t.Fatalf("Apply(Riichi) failed: %v", err)
			}
			if err := s.Apply(event.NewDiscard(actor, tile.MustTileFromCode("W"), false)); err != nil {
				t.Fatalf("Apply(Discard) failed: %v", err)
			}
			if err := s.Apply(event.NewRiichiAccepted(actor, tt.deltas, tt.scores)); err != nil {
				t.Fatalf("Apply(RiichiAccepted) failed: %v", err)
			}

			if got := s.Player(actor).RiichiState(); got != player.RiichiAccepted {
				t.Errorf("RiichiState() = %v, want %v", got, player.RiichiAccepted)
			}
			if got := s.RiichiDeposit(); got != 1 {
				t.Errorf("RiichiDeposit() = %d, want 1", got)
			}
			if got := s.Scores(); got != tt.wantScores {
				t.Errorf("Scores() = %v, want %v", got, tt.wantScores)
			}
		})
	}
}

func TestState_Apply_RiichiAccepted_ReturnsErrorWhenActorIsNotPendingRiichiAcceptancePlayer(t *testing.T) {
	hands := newValidHands()
	hands[0] = riichiReadyHandForTest()
	s := mustNewRoundStateForTest(t, hands)
	riichiActor := seat.MustSeat(0)
	wrongActor := seat.MustSeat(1)

	if err := s.Apply(event.NewDraw(riichiActor, tile.MustTileFromCode("S"))); err != nil {
		t.Fatalf("Apply(Draw) failed: %v", err)
	}
	if err := s.Apply(event.NewRiichi(riichiActor)); err != nil {
		t.Fatalf("Apply(Riichi) failed: %v", err)
	}
	if err := s.Apply(event.NewDiscard(riichiActor, tile.MustTileFromCode("W"), false)); err != nil {
		t.Fatalf("Apply(Discard) failed: %v", err)
	}

	if err := s.Apply(event.NewRiichiAccepted(wrongActor, nil, nil)); err == nil {
		t.Fatal("Apply(RiichiAccepted) succeeded unexpectedly")
	}

	if got := s.Player(riichiActor).RiichiState(); got != player.RiichiDeclared {
		t.Fatalf("RiichiState() = %v, want %v", got, player.RiichiDeclared)
	}
	if got := s.RiichiDeposit(); got != 0 {
		t.Fatalf("RiichiDeposit() = %d, want 0", got)
	}
}

func TestState_Apply_RiichiAccepted_BeforeDeclarationTileCalled(t *testing.T) {
	hands := newValidHands()
	hands[0] = riichiReadyHandForTest()
	hands[3][0] = tile.MustTileFromCode("W")
	hands[3][1] = tile.MustTileFromCode("W")
	s := mustNewRoundStateForTest(t, hands)
	riichiActor := seat.MustSeat(0)
	callActor := seat.MustSeat(3)
	declarationTile := tile.MustTileFromCode("W")

	if err := s.Apply(event.NewDraw(riichiActor, tile.MustTileFromCode("S"))); err != nil {
		t.Fatalf("Apply(Draw) failed: %v", err)
	}
	if err := s.Apply(event.NewRiichi(riichiActor)); err != nil {
		t.Fatalf("Apply(Riichi) failed: %v", err)
	}
	if err := s.Apply(event.NewDiscard(riichiActor, declarationTile, false)); err != nil {
		t.Fatalf("Apply(Discard) failed: %v", err)
	}
	if err := s.Apply(event.NewRiichiAccepted(riichiActor, nil, nil)); err != nil {
		t.Fatalf("Apply(RiichiAccepted) failed: %v", err)
	}
	if err := s.Apply(event.NewPon(callActor, riichiActor, declarationTile, [2]tile.Tile{declarationTile, declarationTile})); err != nil {
		t.Fatalf("Apply(Pon) failed: %v", err)
	}

	if got := s.Player(riichiActor).RiichiState(); got != player.RiichiAccepted {
		t.Fatalf("RiichiState() = %v, want %v", got, player.RiichiAccepted)
	}
	if got := s.RiichiDeposit(); got != 1 {
		t.Fatalf("RiichiDeposit() = %d, want 1", got)
	}
	if got := s.Player(riichiActor).River(); len(got) != 0 {
		t.Fatalf("riichi actor River() = %v, want empty after Pon", got)
	}
	if !s.Player(callActor).CanDiscard() {
		t.Fatal("call actor CanDiscard() = false, want true after Pon")
	}
}
