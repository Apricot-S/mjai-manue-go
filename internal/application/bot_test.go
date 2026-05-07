package application_test

import (
	"errors"
	"testing"

	"github.com/Apricot-S/mjai-manue-go/internal/application"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/action"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/event"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/round"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/seat"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/tile"
)

func TestBot_Process_StartRound(t *testing.T) {
	bot := mustNewBotForTest(t, seat.MustSeat(0))

	got, err := bot.Process(mustNewStartRoundForTest(t, newValidHands()))
	if err != nil {
		t.Fatalf("Process() failed: %v", err)
	}
	if got.Kind() != application.ReactionNone {
		t.Errorf("Kind() = %v, want %v", got.Kind(), application.ReactionNone)
	}
}

func TestBot_Process_DrawSelf(t *testing.T) {
	self := seat.MustSeat(0)
	bot := mustNewBotForTest(t, self)
	if _, err := bot.Process(mustNewStartRoundForTest(t, newValidHands())); err != nil {
		t.Fatalf("Process(StartRound) failed: %v", err)
	}

	drawnTile := tile.MustTileFromCode("6m")
	got, err := bot.Process(event.NewDraw(self, drawnTile))
	if err != nil {
		t.Fatalf("Process(Draw) failed: %v", err)
	}
	if got.Kind() != application.ReactionAction {
		t.Fatalf("Kind() = %v, want %v", got.Kind(), application.ReactionAction)
	}
	discard, ok := got.Action().(*action.Discard)
	if !ok {
		t.Fatalf("Action() = %T, want *action.Discard", got.Action())
	}
	if discard.Actor() != self {
		t.Errorf("Actor() = %v, want %v", discard.Actor(), self)
	}
	if discard.Tile().ID() != drawnTile.ID() {
		t.Errorf("Tile() = %v, want %v", discard.Tile(), drawnTile)
	}
	if !discard.Tsumogiri() {
		t.Error("Tsumogiri() = false, want true")
	}
}

func TestBot_Process_DrawOther(t *testing.T) {
	self := seat.MustSeat(0)
	other := seat.MustSeat(1)
	bot := mustNewBotForTest(t, self)
	if _, err := bot.Process(mustNewStartRoundForTest(t, newValidHands())); err != nil {
		t.Fatalf("Process(StartRound) failed: %v", err)
	}
	drawnTile := tile.MustTileFromCode("6m")
	if _, err := bot.Process(event.NewDraw(self, drawnTile)); err != nil {
		t.Fatalf("Process(self Draw) failed: %v", err)
	}
	if _, err := bot.Process(event.NewDiscard(self, drawnTile, true)); err != nil {
		t.Fatalf("Process(self Discard) failed: %v", err)
	}

	got, err := bot.Process(event.NewDraw(other, tile.MustTileFromCode("6m")))
	if err != nil {
		t.Fatalf("Process(Draw) failed: %v", err)
	}
	if got.Kind() != application.ReactionNone {
		t.Errorf("Kind() = %v, want %v", got.Kind(), application.ReactionNone)
	}
}

func TestBot_Process_Discard(t *testing.T) {
	self := seat.MustSeat(0)
	bot := mustNewBotForTest(t, self)
	if _, err := bot.Process(mustNewStartRoundForTest(t, newValidHands())); err != nil {
		t.Fatalf("Process(StartRound) failed: %v", err)
	}

	drawnTile := tile.MustTileFromCode("6m")
	if _, err := bot.Process(event.NewDraw(self, drawnTile)); err != nil {
		t.Fatalf("Process(Draw) failed: %v", err)
	}
	discard := event.NewDiscard(self, drawnTile, true)

	got, err := bot.Process(discard)
	if err != nil {
		t.Fatalf("Process(Discard) failed: %v", err)
	}
	if got.Kind() != application.ReactionNone {
		t.Errorf("Kind() = %v, want %v", got.Kind(), application.ReactionNone)
	}
}

func TestBot_Process_DrawBeforeStartRound(t *testing.T) {
	bot := mustNewBotForTest(t, seat.MustSeat(0))
	if _, err := bot.Process(event.NewDraw(seat.MustSeat(0), tile.MustTileFromCode("6m"))); err == nil {
		t.Fatal("Process() succeeded unexpectedly")
	}
}

func TestBot_Process_EndRound(t *testing.T) {
	bot := mustNewBotForTest(t, seat.MustSeat(0))
	if _, err := bot.Process(mustNewStartRoundForTest(t, newValidHands())); err != nil {
		t.Fatalf("Process(StartRound) failed: %v", err)
	}

	got, err := bot.Process(event.NewEndRound())
	if err != nil {
		t.Fatalf("Process(EndRound) failed: %v", err)
	}
	if got.Kind() != application.ReactionNone {
		t.Errorf("Kind() = %v, want %v", got.Kind(), application.ReactionNone)
	}
	if _, err := bot.Process(event.NewDraw(seat.MustSeat(0), tile.MustTileFromCode("6m"))); err == nil {
		t.Fatal("Process(Draw) after EndRound succeeded unexpectedly")
	}
}

func TestBot_Process_ReportsRoundStateAfterStateUpdate(t *testing.T) {
	self := seat.MustSeat(0)
	reporter := &recordingRoundStateReporter{}
	bot := application.NewBot(self, newTsumogiriAgentForTest(), reporter)

	if _, err := bot.Process(mustNewStartRoundForTest(t, newValidHands())); err != nil {
		t.Fatalf("Process(StartRound) failed: %v", err)
	}
	if _, err := bot.Process(event.NewDraw(self, tile.MustTileFromCode("6m"))); err != nil {
		t.Fatalf("Process(Draw) failed: %v", err)
	}

	if reporter.calls != 2 {
		t.Errorf("reporter calls = %d, want 2", reporter.calls)
	}
	if reporter.lastBoard == "" {
		t.Error("reported board is empty")
	}
}

func TestBot_Process_ReportsNoRoundStateWhenApplyFails(t *testing.T) {
	self := seat.MustSeat(0)
	reporter := &recordingRoundStateReporter{}
	bot := application.NewBot(self, newTsumogiriAgentForTest(), reporter)

	if _, err := bot.Process(mustNewStartRoundForTest(t, newValidHands())); err != nil {
		t.Fatalf("Process(StartRound) failed: %v", err)
	}
	reporter.calls = 0

	if _, err := bot.Process(event.NewDraw(self, tile.MustTileFromCode("6m"))); err != nil {
		t.Fatalf("Process(first Draw) failed: %v", err)
	}
	reporter.calls = 0
	if _, err := bot.Process(event.NewDraw(self, tile.MustTileFromCode("7m"))); err == nil {
		t.Fatal("Process(second Draw) succeeded unexpectedly")
	}

	if reporter.calls != 0 {
		t.Errorf("reporter calls = %d, want 0", reporter.calls)
	}
}

func TestBot_Process_ReturnsReporterError(t *testing.T) {
	wantErr := errors.New("report failed")
	bot := application.NewBot(
		seat.MustSeat(0),
		newTsumogiriAgentForTest(),
		errorRoundStateReporter{err: wantErr},
	)

	if _, err := bot.Process(mustNewStartRoundForTest(t, newValidHands())); !errors.Is(err, wantErr) {
		t.Errorf("Process() error = %v, want %v", err, wantErr)
	}
}

type recordingRoundStateReporter struct {
	calls     int
	lastBoard string
}

func (r *recordingRoundStateReporter) ReportRoundState(state round.BoardRenderer) error {
	r.calls++
	r.lastBoard = state.RenderBoard()
	return nil
}

type errorRoundStateReporter struct {
	err error
}

func (r errorRoundStateReporter) ReportRoundState(round.BoardRenderer) error {
	return r.err
}
