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

func mustNewRoundStateForTest(t *testing.T, hands [common.NumPlayers][common.InitHandSize]tile.Tile) *State {
	t.Helper()

	validDealer := *seat.MustSeat(0)
	validDora := *tile.MustTileFromCode("E")
	validScores := &[common.NumPlayers]int{25000, 25000, 25000, 25000}

	ev, err := event.NewStartRound(
		wind.East,
		1,
		0,
		0,
		validDealer,
		validDora,
		validScores,
		hands,
	)
	if err != nil {
		t.Fatalf("event.NewStartRound() failed: %v", err)
	}

	s, err := NewState(ev, *validScores)
	if err != nil {
		t.Fatalf("round.NewState() failed: %v", err)
	}
	return s
}

func TestState_Apply_Draw(t *testing.T) {
	t.Run("visible success", func(t *testing.T) {
		s := mustNewRoundStateForTest(t, newValidHands())
		actor := *seat.MustSeat(0)
		drawnTile := *tile.MustTileFromCode("6m")

		before := s.NumLeftTiles()
		if err := s.Apply(event.NewDraw(actor, drawnTile)); err != nil {
			t.Fatalf("Apply() failed: %v", err)
		}

		if got := s.NumLeftTiles(); got != before-1 {
			t.Fatalf("NumLeftTiles() = %d, want %d", got, before-1)
		}

		if got := s.Player(actor).DrawnTile(); got == nil || got.ID() != drawnTile.ID() {
			t.Fatalf("DrawnTile() = %v, want %v", got, drawnTile)
		}
		if !s.Player(actor).CanDiscard() {
			t.Fatalf("actor must be discardable after Draw; CanDiscard() returned false")
		}

		for i := 1; i < common.NumPlayers; i++ {
			playerSeat := *seat.MustSeat(i)
			if got := s.Player(playerSeat).DrawnTile(); got != nil {
				t.Fatalf("player %d DrawnTile() = %v, want nil", i, got)
			}
		}
	})

	t.Run("invisible success (unknown tile allowed)", func(t *testing.T) {
		hands := [4][13]tile.Tile{}
		for p := range common.NumPlayers {
			for i := range common.InitHandSize {
				hands[p][i] = *tile.MustTileFromCode("?")
			}
		}
		s := mustNewRoundStateForTest(t, hands)

		actor := *seat.MustSeat(2)
		drawnTile := *tile.MustTileFromCode("?")

		before := s.NumLeftTiles()
		if err := s.Apply(event.NewDraw(actor, drawnTile)); err != nil {
			t.Fatalf("Apply() failed: %v", err)
		}
		if got := s.NumLeftTiles(); got != before-1 {
			t.Fatalf("NumLeftTiles() = %d, want %d", got, before-1)
		}
		if got := s.Player(actor).DrawnTile(); got == nil || got.ID() != drawnTile.ID() {
			t.Fatalf("DrawnTile() = %v, want %v", got, drawnTile)
		}
		if !s.Player(actor).CanDiscard() {
			t.Fatalf("actor must be discardable after Draw; CanDiscard() returned false")
		}
	})

	t.Run("visible failure (unknown tile)", func(t *testing.T) {
		s := mustNewRoundStateForTest(t, newValidHands())
		actor := *seat.MustSeat(0)
		drawnTile := *tile.MustTileFromCode("?")

		before := s.NumLeftTiles()
		if err := s.Apply(event.NewDraw(actor, drawnTile)); err == nil {
			t.Fatal("Apply() succeeded unexpectedly")
		}
		if got := s.NumLeftTiles(); got != before {
			t.Fatalf("NumLeftTiles() = %d, want %d", got, before)
		}
		if got := s.Player(actor).DrawnTile(); got != nil {
			t.Fatalf("DrawnTile() = %v, want nil", got)
		}
	})

	t.Run("failure does not partially apply (already discardable)", func(t *testing.T) {
		handTiles := newValidHands()[0]
		actorPlayer, err := player.NewVisiblePlayer(handTiles)
		if err != nil {
			t.Fatalf("player.NewVisiblePlayer() failed: %v", err)
		}
		initialDrawnTile := *tile.MustTileFromCode("1m")
		if err := actorPlayer.Draw(initialDrawnTile); err != nil {
			t.Fatalf("actorPlayer.Draw() failed: %v", err)
		}

		players := [common.NumPlayers]player.Player{
			actorPlayer,
			player.NewInvisiblePlayer(),
			player.NewInvisiblePlayer(),
			player.NewInvisiblePlayer(),
		}

		s := NewStateForTest(
			wind.East,
			1,
			0,
			0,
			[common.NumPlayers]int{25000, 25000, 25000, 25000},
			*seat.MustSeat(1),
			*seat.MustSeat(0),
			tile.Tiles{*tile.MustTileFromCode("1m")},
			10,
			players,
		)

		actor := *seat.MustSeat(0)
		anotherTile := *tile.MustTileFromCode("2m")

		before := s.NumLeftTiles()
		if err := s.Apply(event.NewDraw(actor, anotherTile)); err == nil {
			t.Fatal("Apply() succeeded unexpectedly")
		}
		if got := s.NumLeftTiles(); got != before {
			t.Fatalf("NumLeftTiles() = %d, want %d", got, before)
		}
		if got := s.Player(actor).DrawnTile(); got == nil || got.ID() != initialDrawnTile.ID() {
			t.Fatalf("DrawnTile() = %v, want %v", got, initialDrawnTile)
		}
	})

	t.Run("failure (no tiles left)", func(t *testing.T) {
		handTiles := newValidHands()[0]
		actorPlayer, err := player.NewVisiblePlayer(handTiles)
		if err != nil {
			t.Fatalf("player.NewVisiblePlayer() failed: %v", err)
		}

		players := [common.NumPlayers]player.Player{
			actorPlayer,
			player.NewInvisiblePlayer(),
			player.NewInvisiblePlayer(),
			player.NewInvisiblePlayer(),
		}

		s := NewStateForTest(
			wind.East,
			1,
			0,
			0,
			[common.NumPlayers]int{25000, 25000, 25000, 25000},
			*seat.MustSeat(1),
			*seat.MustSeat(0),
			tile.Tiles{*tile.MustTileFromCode("1m")},
			0,
			players,
		)

		actor := *seat.MustSeat(0)
		drawnTile := *tile.MustTileFromCode("6m")

		if err := s.Apply(event.NewDraw(actor, drawnTile)); err == nil {
			t.Fatal("Apply() succeeded unexpectedly")
		}
		if got := s.NumLeftTiles(); got != 0 {
			t.Fatalf("NumLeftTiles() = %d, want %d", got, 0)
		}
		if got := s.Player(actor).DrawnTile(); got != nil {
			t.Fatalf("DrawnTile() = %v, want nil", got)
		}
	})
}
