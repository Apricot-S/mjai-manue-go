package round

import (
	"testing"

	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/common"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/event"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/round/player"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/seat"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/tile"
)

func TestState_Apply_Win(t *testing.T) {
	tests := []struct {
		name       string
		deltas     *[common.NumPlayers]int
		scores     *[common.NumPlayers]int
		wantScores [common.NumPlayers]int
	}{
		{
			name:       "scores",
			scores:     &[common.NumPlayers]int{73000, 9000, 9000, 9000},
			wantScores: [common.NumPlayers]int{73000, 9000, 9000, 9000},
		},
		{
			name:       "deltas",
			deltas:     &[common.NumPlayers]int{48000, -16000, -16000, -16000},
			wantScores: [common.NumPlayers]int{73000, 9000, 9000, 9000},
		},
		{
			name:       "scores take precedence over deltas",
			deltas:     &[common.NumPlayers]int{1, 2, 3, 4},
			scores:     &[common.NumPlayers]int{73000, 9000, 9000, 9000},
			wantScores: [common.NumPlayers]int{73000, 9000, 9000, 9000},
		},
		{
			name:       "no scores or deltas",
			wantScores: [common.NumPlayers]int{25000, 25000, 25000, 25000},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := mustNewRoundStateForTest(t, newValidHands())
			actor := seat.MustSeat(0)
			winningTile := tile.MustTileFromCode("6m")

			if err := s.Apply(event.NewDraw(actor, winningTile)); err != nil {
				t.Fatalf("Apply(Draw) failed: %v", err)
			}
			if err := s.Apply(event.NewWin(
				actor,
				actor,
				&winningTile,
				48000,
				tt.deltas,
				tt.scores,
			)); err != nil {
				t.Fatalf("Apply(Win) failed: %v", err)
			}

			if got := s.Scores(); got != tt.wantScores {
				t.Errorf("Scores() = %v, want %v", got, tt.wantScores)
			}
		})
	}
}

func TestState_Apply_Win_ReturnsErrorBeforeFirstDraw(t *testing.T) {
	s := mustNewRoundStateForTest(t, newValidHands())
	scores := [common.NumPlayers]int{25000, 30800, 34700, 9500}

	if err := s.Apply(event.NewWin(
		seat.MustSeat(2),
		seat.MustSeat(3),
		new(tile.MustTileFromCode("9m")),
		8000,
		nil,
		&scores,
	)); err == nil {
		t.Fatal("Apply(Win) succeeded unexpectedly")
	}

	if got := s.Scores(); got != [common.NumPlayers]int{25000, 25000, 25000, 25000} {
		t.Errorf("Scores() = %v, want unchanged initial scores", got)
	}
}

func TestState_Apply_Win_Renhou(t *testing.T) {
	hands := newValidHands()
	actor := seat.MustSeat(1)
	hands[actor.Index()] = tenpaiHandWaiting36mForTest()
	s := mustNewRoundStateForTest(t, hands)
	target := seat.MustSeat(0)
	winningTile := tile.MustTileFromCode("6m")
	scores := [common.NumPlayers]int{57000, -7000, 25000, 25000}

	if err := s.Apply(event.NewDraw(target, winningTile)); err != nil {
		t.Fatalf("Apply(Draw) failed: %v", err)
	}
	if err := s.Apply(event.NewDiscard(target, winningTile, true)); err != nil {
		t.Fatalf("Apply(Discard) failed: %v", err)
	}
	if err := s.Apply(event.NewWin(
		actor,
		target,
		&winningTile,
		32000,
		nil,
		&scores,
	)); err != nil {
		t.Fatalf("Apply(Win) failed: %v", err)
	}

	if got := s.Scores(); got != scores {
		t.Errorf("Scores() = %v, want %v", got, scores)
	}
}

func TestState_Apply_Win_ReturnsErrorForRonOnOwnDiscardedTile(t *testing.T) {
	hands := newValidHands()
	actor := seat.MustSeat(0)
	hands[actor.Index()] = tenpaiHandWaiting36mForTest()
	s := mustNewRoundStateForTest(t, hands)
	target := seat.MustSeat(3)
	winningTile := tile.MustTileFromCode("6m")
	scores := [common.NumPlayers]int{57000, 25000, 25000, -7000}

	if err := s.Apply(event.NewDraw(actor, winningTile)); err != nil {
		t.Fatalf("Apply(Draw) failed: %v", err)
	}
	if err := s.Apply(event.NewDiscard(actor, winningTile, true)); err != nil {
		t.Fatalf("Apply(Discard) failed: %v", err)
	}
	if err := s.Apply(event.NewDraw(seat.MustSeat(1), tile.MustTileFromCode("7m"))); err != nil {
		t.Fatalf("Apply(Draw) failed: %v", err)
	}
	if err := s.Apply(event.NewDiscard(seat.MustSeat(1), tile.MustTileFromCode("7m"), true)); err != nil {
		t.Fatalf("Apply(Discard) failed: %v", err)
	}
	if err := s.Apply(event.NewDraw(seat.MustSeat(2), tile.MustTileFromCode("8m"))); err != nil {
		t.Fatalf("Apply(Draw) failed: %v", err)
	}
	if err := s.Apply(event.NewDiscard(seat.MustSeat(2), tile.MustTileFromCode("8m"), true)); err != nil {
		t.Fatalf("Apply(Discard) failed: %v", err)
	}
	if err := s.Apply(event.NewDraw(target, winningTile)); err != nil {
		t.Fatalf("Apply(Draw) failed: %v", err)
	}
	if err := s.Apply(event.NewDiscard(target, winningTile, true)); err != nil {
		t.Fatalf("Apply(Discard) failed: %v", err)
	}
	if err := s.Apply(event.NewWin(
		actor,
		target,
		&winningTile,
		32000,
		nil,
		&scores,
	)); err == nil {
		t.Fatal("Apply(Win) succeeded unexpectedly")
	}

	if got := s.Scores(); got != [common.NumPlayers]int{25000, 25000, 25000, 25000} {
		t.Errorf("Scores() = %v, want unchanged initial scores", got)
	}
}

func TestState_Apply_Win_ReturnsErrorForRonOnExtraSafeTile(t *testing.T) {
	hands := newValidHands()
	actor := seat.MustSeat(0)
	hands[actor.Index()] = tenpaiHandWaiting36mForTest()
	s := mustNewRoundStateForTest(t, hands)
	target := seat.MustSeat(3)
	winningTile := tile.MustTileFromCode("6m")
	scores := [common.NumPlayers]int{57000, 25000, 25000, -7000}

	if err := s.Apply(event.NewDraw(actor, tile.MustTileFromCode("7m"))); err != nil {
		t.Fatalf("Apply(Draw) failed: %v", err)
	}
	if err := s.Apply(event.NewDiscard(actor, tile.MustTileFromCode("7m"), true)); err != nil {
		t.Fatalf("Apply(Discard) failed: %v", err)
	}
	if err := s.Apply(event.NewDraw(seat.MustSeat(1), winningTile)); err != nil {
		t.Fatalf("Apply(Draw) failed: %v", err)
	}
	if err := s.Apply(event.NewDiscard(seat.MustSeat(1), winningTile, true)); err != nil {
		t.Fatalf("Apply(Discard) failed: %v", err)
	}
	if err := s.Apply(event.NewDraw(seat.MustSeat(2), tile.MustTileFromCode("7m"))); err != nil {
		t.Fatalf("Apply(Draw) failed: %v", err)
	}
	if err := s.Apply(event.NewDiscard(seat.MustSeat(2), tile.MustTileFromCode("7m"), true)); err != nil {
		t.Fatalf("Apply(Discard) failed: %v", err)
	}
	if err := s.Apply(event.NewDraw(target, winningTile)); err != nil {
		t.Fatalf("Apply(Draw) failed: %v", err)
	}
	if err := s.Apply(event.NewDiscard(target, winningTile, true)); err != nil {
		t.Fatalf("Apply(Discard) failed: %v", err)
	}
	if err := s.Apply(event.NewWin(
		actor,
		target,
		&winningTile,
		32000,
		nil,
		&scores,
	)); err == nil {
		t.Fatal("Apply(Win) succeeded unexpectedly")
	}

	if got := s.Scores(); got != [common.NumPlayers]int{25000, 25000, 25000, 25000} {
		t.Errorf("Scores() = %v, want unchanged initial scores", got)
	}
}

func TestState_Apply_Win_ReturnsErrorForInvisibleRonOnOwnDiscardedTile(t *testing.T) {
	hands := newValidHands()
	actor := seat.MustSeat(0)
	for i := range common.InitHandSize {
		hands[actor.Index()][i] = tile.MustTileFromCode("?")
	}
	s := mustNewRoundStateForTest(t, hands)
	target := seat.MustSeat(3)
	winningTile := tile.MustTileFromCode("6m")
	scores := [common.NumPlayers]int{57000, 25000, 25000, -7000}

	if err := s.Apply(event.NewDraw(actor, tile.MustTileFromCode("?"))); err != nil {
		t.Fatalf("Apply(Draw) failed: %v", err)
	}
	if err := s.Apply(event.NewDiscard(actor, winningTile, false)); err != nil {
		t.Fatalf("Apply(Discard) failed: %v", err)
	}
	if err := s.Apply(event.NewDraw(seat.MustSeat(1), tile.MustTileFromCode("7m"))); err != nil {
		t.Fatalf("Apply(Draw) failed: %v", err)
	}
	if err := s.Apply(event.NewDiscard(seat.MustSeat(1), tile.MustTileFromCode("7m"), true)); err != nil {
		t.Fatalf("Apply(Discard) failed: %v", err)
	}
	if err := s.Apply(event.NewDraw(seat.MustSeat(2), tile.MustTileFromCode("8m"))); err != nil {
		t.Fatalf("Apply(Draw) failed: %v", err)
	}
	if err := s.Apply(event.NewDiscard(seat.MustSeat(2), tile.MustTileFromCode("8m"), true)); err != nil {
		t.Fatalf("Apply(Discard) failed: %v", err)
	}
	if err := s.Apply(event.NewDraw(target, winningTile)); err != nil {
		t.Fatalf("Apply(Draw) failed: %v", err)
	}
	if err := s.Apply(event.NewDiscard(target, winningTile, true)); err != nil {
		t.Fatalf("Apply(Discard) failed: %v", err)
	}
	if err := s.Apply(event.NewWin(
		actor,
		target,
		&winningTile,
		32000,
		nil,
		&scores,
	)); err == nil {
		t.Fatal("Apply(Win) succeeded unexpectedly")
	}

	if got := s.Scores(); got != [common.NumPlayers]int{25000, 25000, 25000, 25000} {
		t.Errorf("Scores() = %v, want unchanged initial scores", got)
	}
}

func TestState_Apply_Win_ReturnsErrorForInvisibleRonOnExtraSafeTile(t *testing.T) {
	hands := newValidHands()
	actor := seat.MustSeat(0)
	for i := range common.InitHandSize {
		hands[actor.Index()][i] = tile.MustTileFromCode("?")
	}
	s := mustNewRoundStateForTest(t, hands)
	target := seat.MustSeat(3)
	winningTile := tile.MustTileFromCode("6m")
	scores := [common.NumPlayers]int{57000, 25000, 25000, -7000}

	if err := s.Apply(event.NewDraw(actor, tile.MustTileFromCode("?"))); err != nil {
		t.Fatalf("Apply(Draw) failed: %v", err)
	}
	if err := s.Apply(event.NewDiscard(actor, tile.MustTileFromCode("7m"), false)); err != nil {
		t.Fatalf("Apply(Discard) failed: %v", err)
	}
	if err := s.Apply(event.NewDraw(seat.MustSeat(1), winningTile)); err != nil {
		t.Fatalf("Apply(Draw) failed: %v", err)
	}
	if err := s.Apply(event.NewDiscard(seat.MustSeat(1), winningTile, true)); err != nil {
		t.Fatalf("Apply(Discard) failed: %v", err)
	}
	if err := s.Apply(event.NewDraw(seat.MustSeat(2), tile.MustTileFromCode("8m"))); err != nil {
		t.Fatalf("Apply(Draw) failed: %v", err)
	}
	if err := s.Apply(event.NewDiscard(seat.MustSeat(2), tile.MustTileFromCode("8m"), true)); err != nil {
		t.Fatalf("Apply(Discard) failed: %v", err)
	}
	if err := s.Apply(event.NewDraw(target, winningTile)); err != nil {
		t.Fatalf("Apply(Draw) failed: %v", err)
	}
	if err := s.Apply(event.NewDiscard(target, winningTile, true)); err != nil {
		t.Fatalf("Apply(Discard) failed: %v", err)
	}
	if err := s.Apply(event.NewWin(
		actor,
		target,
		&winningTile,
		32000,
		nil,
		&scores,
	)); err == nil {
		t.Fatal("Apply(Win) succeeded unexpectedly")
	}

	if got := s.Scores(); got != [common.NumPlayers]int{25000, 25000, 25000, 25000} {
		t.Errorf("Scores() = %v, want unchanged initial scores", got)
	}
}

func TestState_Apply_Win_ReturnsErrorForRonOnSameSymbolRedFiveFuriten(t *testing.T) {
	hands := newValidHands()
	actor := seat.MustSeat(0)
	hands[actor.Index()] = tenpaiHandWaiting58mForTest()
	s := mustNewRoundStateForTest(t, hands)
	target := seat.MustSeat(3)
	safeTile := tile.MustTileFromCode("5m")
	winningTile := tile.MustTileFromCode("5mr")
	scores := [common.NumPlayers]int{57000, 25000, 25000, -7000}

	if err := s.Apply(event.NewDraw(actor, tile.MustTileFromCode("6m"))); err != nil {
		t.Fatalf("Apply(Draw) failed: %v", err)
	}
	if err := s.Apply(event.NewDiscard(actor, tile.MustTileFromCode("6m"), true)); err != nil {
		t.Fatalf("Apply(Discard) failed: %v", err)
	}
	if err := s.Apply(event.NewDraw(seat.MustSeat(1), safeTile)); err != nil {
		t.Fatalf("Apply(Draw) failed: %v", err)
	}
	if err := s.Apply(event.NewDiscard(seat.MustSeat(1), safeTile, true)); err != nil {
		t.Fatalf("Apply(Discard) failed: %v", err)
	}
	if err := s.Apply(event.NewDraw(seat.MustSeat(2), tile.MustTileFromCode("8m"))); err != nil {
		t.Fatalf("Apply(Draw) failed: %v", err)
	}
	if err := s.Apply(event.NewDiscard(seat.MustSeat(2), tile.MustTileFromCode("8m"), true)); err != nil {
		t.Fatalf("Apply(Discard) failed: %v", err)
	}
	if err := s.Apply(event.NewDraw(target, winningTile)); err != nil {
		t.Fatalf("Apply(Draw) failed: %v", err)
	}
	if err := s.Apply(event.NewDiscard(target, winningTile, true)); err != nil {
		t.Fatalf("Apply(Discard) failed: %v", err)
	}
	if err := s.Apply(event.NewWin(
		actor,
		target,
		&winningTile,
		32000,
		nil,
		&scores,
	)); err == nil {
		t.Fatal("Apply(Win) succeeded unexpectedly")
	}

	if got := s.Scores(); got != [common.NumPlayers]int{25000, 25000, 25000, 25000} {
		t.Errorf("Scores() = %v, want unchanged initial scores", got)
	}
}

func TestState_Apply_Win_RobbingKan(t *testing.T) {
	s := newStateBeforePromotedKanForTest(t, 10, 0)
	actor := seat.MustSeat(1)
	target := seat.MustSeat(3)
	added := tile.MustTileFromCode("E")
	scores := [common.NumPlayers]int{25000, 57000, 25000, -7000}

	if err := s.Apply(event.NewPromotedKan(target, added, [3]tile.Tile{added, added, added})); err != nil {
		t.Fatalf("Apply(PromotedKan) failed: %v", err)
	}
	if err := s.Apply(event.NewWin(
		actor,
		target,
		&added,
		32000,
		nil,
		&scores,
	)); err != nil {
		t.Fatalf("Apply(Win) failed: %v", err)
	}

	if got := s.Scores(); got != scores {
		t.Errorf("Scores() = %v, want %v", got, scores)
	}
}

func TestState_Apply_Win_ReturnsErrorForFuritenRobbingKan(t *testing.T) {
	s := newStateBeforePromotedKanForTest(t, 10, 0)
	actor := seat.MustSeat(1)
	target := seat.MustSeat(3)
	added := tile.MustTileFromCode("E")
	scores := [common.NumPlayers]int{25000, 57000, 25000, -7000}
	actorPlayer, err := player.NewVisiblePlayer(tenpaiHandWaitingEastForTest())
	if err != nil {
		t.Fatalf("player.NewVisiblePlayer() failed: %v", err)
	}
	actorPlayer.AddExtraSafeTiles(added)
	s.players[actor.Index()] = actorPlayer

	if err := s.Apply(event.NewPromotedKan(target, added, [3]tile.Tile{added, added, added})); err != nil {
		t.Fatalf("Apply(PromotedKan) failed: %v", err)
	}
	if err := s.Apply(event.NewWin(
		actor,
		target,
		&added,
		32000,
		nil,
		&scores,
	)); err == nil {
		t.Fatal("Apply(Win) succeeded unexpectedly")
	}

	if got := s.Scores(); got != [common.NumPlayers]int{25000, 25000, 25000, 25000} {
		t.Errorf("Scores() = %v, want unchanged initial scores", got)
	}
}

func TestState_Apply_Win_ReturnsErrorDuringRobbingKanWithDifferentWinningTile(t *testing.T) {
	s := newStateBeforePromotedKanForTest(t, 10, 0)
	actor := seat.MustSeat(1)
	target := seat.MustSeat(3)
	added := tile.MustTileFromCode("E")
	winningTile := tile.MustTileFromCode("S")
	scores := [common.NumPlayers]int{25000, 57000, 25000, -7000}

	if got := s.Player(target).River(); len(got) != 1 || got[0] != winningTile {
		t.Fatalf("target River() = %v, want [%v]", got, winningTile)
	}
	if err := s.Apply(event.NewPromotedKan(target, added, [3]tile.Tile{added, added, added})); err != nil {
		t.Fatalf("Apply(PromotedKan) failed: %v", err)
	}
	if err := s.Apply(event.NewWin(
		actor,
		target,
		&winningTile,
		32000,
		nil,
		&scores,
	)); err == nil {
		t.Fatal("Apply(Win) succeeded unexpectedly")
	}

	if got := s.Scores(); got != [common.NumPlayers]int{25000, 25000, 25000, 25000} {
		t.Errorf("Scores() = %v, want unchanged initial scores", got)
	}
}

func TestState_Apply_Win_ReturnsErrorForRonOnMissedPromotedKanTile(t *testing.T) {
	s := newStateBeforePromotedKanForTest(t, 10, 0)
	actor := seat.MustSeat(1)
	target := seat.MustSeat(3)
	added := tile.MustTileFromCode("E")
	scores := [common.NumPlayers]int{25000, 57000, 25000, -7000}
	actorPlayer, err := player.NewVisiblePlayer(tenpaiHandWaitingEastForTest())
	if err != nil {
		t.Fatalf("player.NewVisiblePlayer() failed: %v", err)
	}
	s.players[actor.Index()] = actorPlayer

	if err := s.Apply(event.NewPromotedKan(target, added, [3]tile.Tile{added, added, added})); err != nil {
		t.Fatalf("Apply(PromotedKan) failed: %v", err)
	}
	if err := s.Apply(event.NewDraw(target, added)); err != nil {
		t.Fatalf("Apply(Draw) failed: %v", err)
	}
	if err := s.Apply(event.NewDora(tile.MustTileFromCode("6p"))); err != nil {
		t.Fatalf("Apply(Dora) failed: %v", err)
	}
	if err := s.Apply(event.NewDiscard(target, added, true)); err != nil {
		t.Fatalf("Apply(Discard) failed: %v", err)
	}
	if err := s.Apply(event.NewWin(
		actor,
		target,
		&added,
		32000,
		nil,
		&scores,
	)); err == nil {
		t.Fatal("Apply(Win) succeeded unexpectedly")
	}

	if got := s.Scores(); got != [common.NumPlayers]int{25000, 25000, 25000, 25000} {
		t.Errorf("Scores() = %v, want unchanged initial scores", got)
	}
}

func tenpaiHandWaiting36mForTest() [common.InitHandSize]tile.Tile {
	return handTilesForTest(
		"1m", "2m", "3m",
		"1p", "2p", "3p",
		"1s", "2s", "3s",
		"4m", "5m",
		"5p", "5p",
	)
}

func tenpaiHandWaiting58mForTest() [common.InitHandSize]tile.Tile {
	return handTilesForTest(
		"1m", "2m", "3m",
		"1p", "2p", "3p",
		"1s", "2s", "3s",
		"6m", "7m",
		"5p", "5p",
	)
}

func tenpaiHandWaitingEastForTest() [common.InitHandSize]tile.Tile {
	return handTilesForTest(
		"1m", "1m", "1m",
		"2p", "2p", "2p",
		"3s", "3s", "3s",
		"S", "S", "S",
		"E",
	)
}

func handTilesForTest(codes ...string) [common.InitHandSize]tile.Tile {
	var handTiles [common.InitHandSize]tile.Tile
	for i, code := range codes {
		handTiles[i] = tile.MustTileFromCode(code)
	}
	return handTiles
}

func TestState_Apply_Win_TsumoWithoutWinningTile(t *testing.T) {
	s := mustNewRoundStateForTest(t, newValidHands())
	actor := seat.MustSeat(0)
	drawnTile := tile.MustTileFromCode("6m")
	scores := [common.NumPlayers]int{73000, 9000, 9000, 9000}

	if err := s.Apply(event.NewDraw(actor, drawnTile)); err != nil {
		t.Fatalf("Apply(Draw) failed: %v", err)
	}
	if err := s.Apply(event.NewWin(
		actor,
		actor,
		nil,
		48000,
		nil,
		&scores,
	)); err != nil {
		t.Fatalf("Apply(Win) failed: %v", err)
	}

	if got := s.Scores(); got != scores {
		t.Errorf("Scores() = %v, want %v", got, scores)
	}
}

func TestState_Apply_Win_InvisibleTsumo(t *testing.T) {
	hands := [common.NumPlayers][common.InitHandSize]tile.Tile{}
	for p := range common.NumPlayers {
		for i := range common.InitHandSize {
			hands[p][i] = tile.MustTileFromCode("?")
		}
	}
	s := mustNewRoundStateForTest(t, hands)
	actor := seat.MustSeat(0)
	unknownDrawnTile := tile.MustTileFromCode("?")
	winningTile := tile.MustTileFromCode("6m")
	scores := [common.NumPlayers]int{73000, 9000, 9000, 9000}

	if err := s.Apply(event.NewDraw(actor, unknownDrawnTile)); err != nil {
		t.Fatalf("Apply(Draw) failed: %v", err)
	}
	if err := s.Apply(event.NewWin(
		actor,
		actor,
		&winningTile,
		48000,
		nil,
		&scores,
	)); err != nil {
		t.Fatalf("Apply(Win) failed: %v", err)
	}

	if got := s.Scores(); got != scores {
		t.Errorf("Scores() = %v, want %v", got, scores)
	}
}
