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

	ev := event.NewStartRound(
		wind.East,
		1,
		0,
		0,
		validDealer,
		validDora,
		validScores,
		hands,
	)

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

func TestState_Apply_Discard(t *testing.T) {
	t.Run("visible tsumogiri success", func(t *testing.T) {
		s := mustNewRoundStateForTest(t, newValidHands())
		actor := *seat.MustSeat(0)
		discardedTile := *tile.MustTileFromCode("6m")

		if err := s.Apply(event.NewDraw(actor, discardedTile)); err != nil {
			t.Fatalf("Apply(Draw) failed: %v", err)
		}
		before := s.NumLeftTiles()

		ev := event.NewDiscard(actor, discardedTile, true)
		if err := s.Apply(ev); err != nil {
			t.Fatalf("Apply(Discard) failed: %v", err)
		}

		if got := s.NumLeftTiles(); got != before {
			t.Errorf("NumLeftTiles() = %d, want %d", got, before)
		}
		if got := s.Player(actor).DrawnTile(); got != nil {
			t.Fatalf("DrawnTile() = %v, want nil", got)
		}
		if s.Player(actor).CanDiscard() {
			t.Fatal("CanDiscard() = true, want false")
		}
		if got := s.Player(actor).River(); len(got) != 1 || got[0].ID() != discardedTile.ID() {
			t.Fatalf("River() = %v, want [%v]", got, discardedTile)
		}
		if got := s.Player(actor).DiscardedTiles(); len(got) != 1 || got[0].ID() != discardedTile.ID() {
			t.Fatalf("DiscardedTiles() = %v, want [%v]", got, discardedTile)
		}
	})

	t.Run("visible hand discard success", func(t *testing.T) {
		s := mustNewRoundStateForTest(t, newValidHands())
		actor := *seat.MustSeat(0)
		drawnTile := *tile.MustTileFromCode("6m")
		discardedTile := newValidHands()[0][0]

		if err := s.Apply(event.NewDraw(actor, drawnTile)); err != nil {
			t.Fatalf("Apply(Draw) failed: %v", err)
		}
		ev := event.NewDiscard(actor, discardedTile, false)
		if err := s.Apply(ev); err != nil {
			t.Fatalf("Apply(Discard) failed: %v", err)
		}

		if got := s.Player(actor).DrawnTile(); got != nil {
			t.Fatalf("DrawnTile() = %v, want nil", got)
		}
		if got := s.Player(actor).River(); len(got) != 1 || got[0].ID() != discardedTile.ID() {
			t.Fatalf("River() = %v, want [%v]", got, discardedTile)
		}
		handTiles := s.Player(actor).HandTiles()
		foundDrawnTile := false
		foundDiscardedTile := false
		for _, handTile := range handTiles {
			if handTile.ID() == drawnTile.ID() {
				foundDrawnTile = true
			}
			if handTile.ID() == discardedTile.ID() {
				foundDiscardedTile = true
			}
		}
		if !foundDrawnTile {
			t.Fatalf("HandTiles() = %v, want drawn tile %v", handTiles, drawnTile)
		}
		if foundDiscardedTile {
			t.Fatalf("HandTiles() = %v, must not contain discarded tile %v", handTiles, discardedTile)
		}
	})

	t.Run("invisible success", func(t *testing.T) {
		hands := [4][13]tile.Tile{}
		for p := range common.NumPlayers {
			for i := range common.InitHandSize {
				hands[p][i] = *tile.MustTileFromCode("?")
			}
		}
		s := mustNewRoundStateForTest(t, hands)
		actor := *seat.MustSeat(2)
		drawnTile := *tile.MustTileFromCode("?")
		discardedTile := *tile.MustTileFromCode("E")

		if err := s.Apply(event.NewDraw(actor, drawnTile)); err != nil {
			t.Fatalf("Apply(Draw) failed: %v", err)
		}
		ev := event.NewDiscard(actor, discardedTile, true)
		if err := s.Apply(ev); err != nil {
			t.Fatalf("Apply(Discard) failed: %v", err)
		}

		if got := s.Player(actor).DrawnTile(); got != nil {
			t.Fatalf("DrawnTile() = %v, want nil", got)
		}
		if got := s.Player(actor).River(); len(got) != 1 || got[0].ID() != discardedTile.ID() {
			t.Fatalf("River() = %v, want [%v]", got, discardedTile)
		}
	})

	t.Run("failure does not partially apply", func(t *testing.T) {
		s := mustNewRoundStateForTest(t, newValidHands())
		actor := *seat.MustSeat(0)
		drawnTile := *tile.MustTileFromCode("6m")
		wrongDiscardedTile := *tile.MustTileFromCode("7m")

		if err := s.Apply(event.NewDraw(actor, drawnTile)); err != nil {
			t.Fatalf("Apply(Draw) failed: %v", err)
		}
		ev := event.NewDiscard(actor, wrongDiscardedTile, true)
		if err := s.Apply(ev); err == nil {
			t.Fatal("Apply(Discard) succeeded unexpectedly")
		}

		if got := s.Player(actor).DrawnTile(); got == nil || got.ID() != drawnTile.ID() {
			t.Fatalf("DrawnTile() = %v, want %v", got, drawnTile)
		}
		if got := s.Player(actor).River(); len(got) != 0 {
			t.Fatalf("River() = %v, want empty", got)
		}
		if got := s.Player(actor).DiscardedTiles(); len(got) != 0 {
			t.Fatalf("DiscardedTiles() = %v, want empty", got)
		}
	})
}

func TestState_Apply_RiichiAndAccepted(t *testing.T) {
	hands := newValidHands()
	hands[0] = [common.InitHandSize]tile.Tile{
		*tile.MustTileFromCode("1m"), *tile.MustTileFromCode("2m"), *tile.MustTileFromCode("3m"),
		*tile.MustTileFromCode("4p"), *tile.MustTileFromCode("5p"), *tile.MustTileFromCode("6p"),
		*tile.MustTileFromCode("7s"), *tile.MustTileFromCode("8s"), *tile.MustTileFromCode("9s"),
		*tile.MustTileFromCode("E"), *tile.MustTileFromCode("E"), *tile.MustTileFromCode("S"),
		*tile.MustTileFromCode("W"),
	}
	s := mustNewRoundStateForTest(t, hands)
	actor := *seat.MustSeat(0)

	if err := s.Apply(event.NewDraw(actor, *tile.MustTileFromCode("S"))); err != nil {
		t.Fatalf("Apply(Draw) failed: %v", err)
	}
	if err := s.Apply(event.NewRiichi(actor)); err != nil {
		t.Fatalf("Apply(Riichi) failed: %v", err)
	}
	if got := s.Player(actor).RiichiState(); got != player.RiichiDeclared {
		t.Fatalf("RiichiState() = %v, want %v", got, player.RiichiDeclared)
	}

	discard := event.NewDiscard(actor, *tile.MustTileFromCode("W"), false)
	if err := s.Apply(discard); err != nil {
		t.Fatalf("Apply(Discard) failed: %v", err)
	}

	scores := [common.NumPlayers]int{24000, 25000, 25000, 25000}
	deltas := [common.NumPlayers]int{-1000, 0, 0, 0}
	if err := s.Apply(event.NewRiichiAccepted(actor, &deltas, &scores)); err != nil {
		t.Fatalf("Apply(RiichiAccepted) failed: %v", err)
	}
	if got := s.Player(actor).RiichiState(); got != player.RiichiAccepted {
		t.Errorf("RiichiState() = %v, want %v", got, player.RiichiAccepted)
	}
	if got := s.RiichiDeposit(); got != 1 {
		t.Errorf("RiichiDeposit() = %d, want 1", got)
	}
	if got := s.Scores(); got != scores {
		t.Errorf("Scores() = %v, want %v", got, scores)
	}
}

func TestState_Apply_RiichiAcceptedActorOnlySubtractsDeposit(t *testing.T) {
	hands := newValidHands()
	hands[1] = [common.InitHandSize]tile.Tile{
		*tile.MustTileFromCode("1m"), *tile.MustTileFromCode("2m"), *tile.MustTileFromCode("3m"),
		*tile.MustTileFromCode("4p"), *tile.MustTileFromCode("5p"), *tile.MustTileFromCode("6p"),
		*tile.MustTileFromCode("7s"), *tile.MustTileFromCode("8s"), *tile.MustTileFromCode("9s"),
		*tile.MustTileFromCode("E"), *tile.MustTileFromCode("E"), *tile.MustTileFromCode("S"),
		*tile.MustTileFromCode("W"),
	}
	s := mustNewRoundStateForTest(t, hands)
	actor := *seat.MustSeat(1)

	if err := s.Apply(event.NewDraw(actor, *tile.MustTileFromCode("S"))); err != nil {
		t.Fatalf("Apply(Draw) failed: %v", err)
	}
	if err := s.Apply(event.NewRiichi(actor)); err != nil {
		t.Fatalf("Apply(Riichi) failed: %v", err)
	}
	discard := event.NewDiscard(actor, *tile.MustTileFromCode("W"), false)
	if err := s.Apply(discard); err != nil {
		t.Fatalf("Apply(Discard) failed: %v", err)
	}

	before := s.Scores()
	if err := s.Apply(event.NewRiichiAccepted(actor, nil, nil)); err != nil {
		t.Fatalf("Apply(RiichiAccepted) failed: %v", err)
	}
	want := before
	want[actor.Index()] -= 1000
	if got := s.Scores(); got != want {
		t.Errorf("Scores() = %v, want %v", got, want)
	}
	if got := s.RiichiDeposit(); got != 1 {
		t.Errorf("RiichiDeposit() = %d, want 1", got)
	}
}

func TestState_Apply_Pon(t *testing.T) {
	hands := newValidHands()
	hands[0][1] = *tile.MustTileFromCode("1s")
	s := mustNewRoundStateForTest(t, hands)
	actor := *seat.MustSeat(0)
	target := *seat.MustSeat(3)
	taken := *tile.MustTileFromCode("1s")

	if err := s.Apply(event.NewDraw(target, taken)); err != nil {
		t.Fatalf("Apply(Draw) failed: %v", err)
	}
	discard := event.NewDiscard(target, taken, true)
	if err := s.Apply(discard); err != nil {
		t.Fatalf("Apply(Discard) failed: %v", err)
	}

	pon := event.NewPon(actor, target, taken, [2]tile.Tile{taken, taken})
	if err := s.Apply(pon); err != nil {
		t.Fatalf("Apply(Pon) failed: %v", err)
	}

	if got := s.Player(target).River(); len(got) != 0 {
		t.Errorf("target River() = %v, want empty", got)
	}
	if got := s.Player(target).DiscardedTiles(); len(got) != 1 || got[0] != taken {
		t.Errorf("target DiscardedTiles() = %v, want [%v]", got, taken)
	}
	if got := s.Player(actor).Melds(); len(got) != 1 {
		t.Fatalf("actor Melds() length = %d, want 1", len(got))
	}
	if !s.Player(actor).CanDiscard() {
		t.Error("actor CanDiscard() = false, want true after Pon")
	}
}

func TestState_Apply_DoraAndRoundResultScores(t *testing.T) {
	s := mustNewRoundStateForTest(t, newValidHands())
	doraIndicator := *tile.MustTileFromCode("6p")
	dora := event.NewDora(doraIndicator)
	if err := s.Apply(dora); err != nil {
		t.Fatalf("Apply(Dora) failed: %v", err)
	}
	if got := s.DoraIndicators(); len(got) != 2 || got[1] != doraIndicator {
		t.Fatalf("DoraIndicators() = %v, want appended %v", got, doraIndicator)
	}

	if err := s.Apply(event.NewDora(*tile.MustTileFromCode("?"))); err == nil {
		t.Fatal("Apply(Dora unknown) succeeded unexpectedly")
	}

	winScores := [common.NumPlayers]int{25000, 30800, 34700, 9500}
	win := event.NewWin(
		*seat.MustSeat(2),
		*seat.MustSeat(3),
		tile.MustTileFromCode("9m"),
		8000,
		nil,
		&winScores,
	)
	if err := s.Apply(win); err != nil {
		t.Fatalf("Apply(Win) failed: %v", err)
	}
	if got := s.Scores(); got != winScores {
		t.Errorf("Scores() after Win = %v, want %v", got, winScores)
	}

	drawRoundScores := [common.NumPlayers]int{23500, 26500, 23500, 26500}
	drawRound := event.NewDrawRound(
		"fanpai",
		&[common.NumPlayers]bool{false, true, false, true},
		nil,
		&drawRoundScores,
	)
	if err := s.Apply(drawRound); err != nil {
		t.Fatalf("Apply(DrawRound) failed: %v", err)
	}
	if got := s.Scores(); got != drawRoundScores {
		t.Errorf("Scores() after DrawRound = %v, want %v", got, drawRoundScores)
	}
}

func TestState_Apply_RoundResultDeltas(t *testing.T) {
	s := mustNewRoundStateForTest(t, newValidHands())
	before := s.Scores()
	deltas := [common.NumPlayers]int{0, 0, 10300, -8300}

	win := event.NewWin(
		*seat.MustSeat(2),
		*seat.MustSeat(3),
		tile.MustTileFromCode("9m"),
		8000,
		&deltas,
		nil,
	)
	if err := s.Apply(win); err != nil {
		t.Fatalf("Apply(Win) failed: %v", err)
	}
	want := before
	for i, delta := range deltas {
		want[i] += delta
	}
	if got := s.Scores(); got != want {
		t.Errorf("Scores() after Win deltas = %v, want %v", got, want)
	}
}

func TestState_Apply_RoundResultWithoutScoresOrDeltas(t *testing.T) {
	s := mustNewRoundStateForTest(t, newValidHands())
	before := s.Scores()
	drawRound := event.NewDrawRound("", nil, nil, nil)

	if err := s.Apply(drawRound); err != nil {
		t.Fatalf("Apply(DrawRound) failed: %v", err)
	}
	if got := s.Scores(); got != before {
		t.Errorf("Scores() = %v, want unchanged %v", got, before)
	}
}
