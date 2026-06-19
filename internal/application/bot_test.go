package application_test

import (
	"errors"
	"testing"

	"github.com/Apricot-S/mjai-manue-go/internal/application"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/ai"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/action"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/common"
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

func TestBot_Process_ReachAcceptedDoesNotRepeatDeclarationTileCall(t *testing.T) {
	self := seat.MustSeat(3)
	riichiActor := seat.MustSeat(0)
	declarationTile := tile.MustTileFromCode("W")
	hands := newValidHands()
	hands[0] = riichiReadyHandForApplicationTest()
	hands[3] = ponHandForApplicationTest("W", "W")
	bot := application.NewBot(self, firstLegalActionAgent{}, nil)
	if _, err := bot.Process(mustNewStartRoundForTest(t, hands)); err != nil {
		t.Fatalf("Process(StartRound) failed: %v", err)
	}
	if _, err := bot.Process(event.NewDraw(riichiActor, tile.MustTileFromCode("S"))); err != nil {
		t.Fatalf("Process(Draw) failed: %v", err)
	}
	if _, err := bot.Process(event.NewRiichi(riichiActor)); err != nil {
		t.Fatalf("Process(Riichi) failed: %v", err)
	}

	firstReaction, err := bot.Process(event.NewDiscard(riichiActor, declarationTile, false))
	if err != nil {
		t.Fatalf("Process(Discard) failed: %v", err)
	}
	if firstReaction.Kind() != application.ReactionAction {
		t.Fatalf("Kind() = %v, want %v", firstReaction.Kind(), application.ReactionAction)
	}
	if _, ok := firstReaction.Action().(*action.Pon); !ok {
		t.Fatalf("Action() = %T, want *action.Pon", firstReaction.Action())
	}

	got, err := bot.Process(event.NewRiichiAccepted(riichiActor, nil, nil))
	if err != nil {
		t.Fatalf("Process(RiichiAccepted) failed: %v", err)
	}
	if got.Kind() != application.ReactionNone {
		t.Fatalf("Kind() = %v, want %v; action = %T", got.Kind(), application.ReactionNone, got.Action())
	}

	if _, err := bot.Process(event.NewPon(self, riichiActor, declarationTile, [2]tile.Tile{declarationTile, declarationTile})); err != nil {
		t.Fatalf("Process(Pon) failed: %v", err)
	}
	got, err = bot.Process(event.NewDiscard(self, tile.MustTileFromCode("9p"), false))
	if err != nil {
		t.Fatalf("Process(self Discard) failed: %v", err)
	}
	if got.Kind() != application.ReactionNone {
		t.Fatalf("Kind() after self discard = %v, want %v", got.Kind(), application.ReactionNone)
	}
}

func TestBot_Process_CalledKanReplacementDrawKeepsDiscardAcrossDora(t *testing.T) {
	self := seat.MustSeat(1)
	target := seat.MustSeat(0)
	kanTile := tile.MustTileFromCode("E")
	replacementTile := tile.MustTileFromCode("W")
	hands := newValidHands()
	hands[1] = calledKanHandForApplicationTest("E", "E", "E")
	bot := application.NewBot(self, firstLegalActionAgent{}, nil)
	if _, err := bot.Process(mustNewStartRoundForTest(t, hands)); err != nil {
		t.Fatalf("Process(StartRound) failed: %v", err)
	}
	if _, err := bot.Process(event.NewDraw(target, kanTile)); err != nil {
		t.Fatalf("Process(target Draw) failed: %v", err)
	}

	firstReaction, err := bot.Process(event.NewDiscard(target, kanTile, true))
	if err != nil {
		t.Fatalf("Process(target Discard) failed: %v", err)
	}
	if firstReaction.Kind() != application.ReactionAction {
		t.Fatalf("Kind() = %v, want %v", firstReaction.Kind(), application.ReactionAction)
	}

	if _, err := bot.Process(event.NewCalledKan(self, target, kanTile, [3]tile.Tile{kanTile, kanTile, kanTile})); err != nil {
		t.Fatalf("Process(CalledKan) failed: %v", err)
	}
	replacementReaction, err := bot.Process(event.NewDraw(self, replacementTile))
	if err != nil {
		t.Fatalf("Process(replacement Draw) failed: %v", err)
	}
	if replacementReaction.Kind() != application.ReactionAction {
		t.Fatalf("replacement Draw Kind() = %v, want %v", replacementReaction.Kind(), application.ReactionAction)
	}
	discard, ok := replacementReaction.Action().(*action.Discard)
	if !ok {
		t.Fatalf("Action() = %T, want *action.Discard", replacementReaction.Action())
	}
	if discard.Actor() != self {
		t.Fatalf("Discard actor = %v, want %v", discard.Actor(), self)
	}

	got, err := bot.Process(event.NewDora(tile.MustTileFromCode("6p")))
	if err != nil {
		t.Fatalf("Process(Dora) failed: %v", err)
	}
	if got.Kind() != application.ReactionNone {
		t.Fatalf("Kind() = %v, want %v; action = %T", got.Kind(), application.ReactionNone, got.Action())
	}
	if _, err := bot.Process(event.NewDiscard(self, discard.Tile(), discard.Tsumogiri())); err != nil {
		t.Fatalf("Process(replacement Discard) failed: %v", err)
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
	reporter := &recordingReporter{}
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

func TestBot_Process_ReportsDecisionTrace(t *testing.T) {
	self := seat.MustSeat(0)
	reporter := &recordingReporter{}
	bot := application.NewBot(self, traceAgent{trace: "evaluation trace\n"}, reporter)

	if _, err := bot.Process(mustNewStartRoundForTest(t, newValidHands())); err != nil {
		t.Fatalf("Process(StartRound) failed: %v", err)
	}
	if _, err := bot.Process(event.NewDraw(self, tile.MustTileFromCode("6m"))); err != nil {
		t.Fatalf("Process(Draw) failed: %v", err)
	}

	if reporter.lastTrace != "evaluation trace\n" {
		t.Errorf("reported trace = %q, want evaluation trace", reporter.lastTrace)
	}
}

func TestBot_Process_ReportsNoRoundStateWhenApplyFails(t *testing.T) {
	self := seat.MustSeat(0)
	reporter := &recordingReporter{}
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
		errorReporter{err: wantErr},
	)

	if _, err := bot.Process(mustNewStartRoundForTest(t, newValidHands())); !errors.Is(err, wantErr) {
		t.Errorf("Process() error = %v, want %v", err, wantErr)
	}
}

type recordingReporter struct {
	calls     int
	lastBoard string
	lastTrace string
}

func (r *recordingReporter) ReportRoundState(state round.BoardRenderer) error {
	r.calls++
	r.lastBoard = state.RenderBoard()
	return nil
}

func (r *recordingReporter) ReportDecisionTrace(trace string) error {
	r.lastTrace = trace
	return nil
}

type errorReporter struct {
	err error
}

func (r errorReporter) ReportRoundState(round.BoardRenderer) error {
	return r.err
}

func (r errorReporter) ReportDecisionTrace(string) error {
	return r.err
}

type traceAgent struct {
	trace string
}

func (traceAgent) Reset() {}

func (a traceAgent) Decide(request ai.Request) (ai.Decision, error) {
	decision, err := ai.NewTsumogiriAgent().Decide(request)
	if err != nil {
		return ai.Decision{}, err
	}
	decision.Trace = a.trace
	return decision, nil
}

type firstLegalActionAgent struct{}

func (firstLegalActionAgent) Reset() {}

func (firstLegalActionAgent) Decide(request ai.Request) (ai.Decision, error) {
	actions, err := request.Round.LegalActions(request.Self)
	if err != nil {
		return ai.Decision{}, err
	}
	if len(actions) == 0 {
		return ai.Decision{}, nil
	}
	return ai.Decision{Action: actions[0]}, nil
}

func riichiReadyHandForApplicationTest() [common.InitHandSize]tile.Tile {
	return [common.InitHandSize]tile.Tile{
		tile.MustTileFromCode("1m"), tile.MustTileFromCode("2m"), tile.MustTileFromCode("3m"),
		tile.MustTileFromCode("4p"), tile.MustTileFromCode("5p"), tile.MustTileFromCode("6p"),
		tile.MustTileFromCode("7s"), tile.MustTileFromCode("8s"), tile.MustTileFromCode("9s"),
		tile.MustTileFromCode("E"), tile.MustTileFromCode("E"), tile.MustTileFromCode("S"),
		tile.MustTileFromCode("W"),
	}
}

func ponHandForApplicationTest(firstCode, secondCode string) [common.InitHandSize]tile.Tile {
	return [common.InitHandSize]tile.Tile{
		tile.MustTileFromCode(firstCode),
		tile.MustTileFromCode(secondCode),
		tile.MustTileFromCode("1m"),
		tile.MustTileFromCode("2m"),
		tile.MustTileFromCode("3m"),
		tile.MustTileFromCode("1p"),
		tile.MustTileFromCode("2p"),
		tile.MustTileFromCode("3p"),
		tile.MustTileFromCode("1s"),
		tile.MustTileFromCode("2s"),
		tile.MustTileFromCode("3s"),
		tile.MustTileFromCode("9p"),
		tile.MustTileFromCode("9s"),
	}
}

func calledKanHandForApplicationTest(firstCode, secondCode, thirdCode string) [common.InitHandSize]tile.Tile {
	return [common.InitHandSize]tile.Tile{
		tile.MustTileFromCode(firstCode),
		tile.MustTileFromCode(secondCode),
		tile.MustTileFromCode(thirdCode),
		tile.MustTileFromCode("1m"),
		tile.MustTileFromCode("2m"),
		tile.MustTileFromCode("3m"),
		tile.MustTileFromCode("1p"),
		tile.MustTileFromCode("2p"),
		tile.MustTileFromCode("3p"),
		tile.MustTileFromCode("1s"),
		tile.MustTileFromCode("2s"),
		tile.MustTileFromCode("3s"),
		tile.MustTileFromCode("9p"),
	}
}
