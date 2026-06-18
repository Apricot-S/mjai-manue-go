package mjairuntime_test

import (
	"testing"

	"github.com/Apricot-S/mjai-manue-go/internal/adapter/mjai/inbound"
	"github.com/Apricot-S/mjai-manue-go/internal/adapter/mjai/outbound"
	"github.com/Apricot-S/mjai-manue-go/internal/adapter/mjai/runtime"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/ai"
)

func TestDriver_HandleHello(t *testing.T) {
	driver := mjairuntime.NewDriver("tsumogiri", "default", 0, ai.NewTsumogiriAgent(), nil)

	msg, err := driver.Handle(&inbound.Hello{Type: "hello"})
	if err != nil {
		t.Fatalf("Handle() failed: %v", err)
	}
	got, ok := msg.(*outbound.Join)
	if !ok {
		t.Fatalf("Handle() = %T, want *outbound.Join", msg)
	}
	if got.Name != "tsumogiri" {
		t.Errorf("Name = %q, want tsumogiri", got.Name)
	}
	if got.Room != "default" {
		t.Errorf("Room = %q, want default", got.Room)
	}
}

func TestDriver_HandleStartGameCreatesBotWithoutOutput(t *testing.T) {
	driver := mjairuntime.NewDriver("tsumogiri", "default", 0, ai.NewTsumogiriAgent(), nil)

	msg, err := driver.Handle(&inbound.StartGame{Type: "start_game", ID: intPtr(0)})
	if err != nil {
		t.Fatalf("Handle() failed: %v", err)
	}
	if msg != nil {
		t.Errorf("Handle() = %T, want nil", msg)
	}
}

func TestDriver_HandleStartGameResetsAgent(t *testing.T) {
	agent := &recordingAgent{}
	driver := mjairuntime.NewDriver("manue", "default", 0, agent, nil)

	if _, err := driver.Handle(&inbound.StartGame{Type: "start_game", ID: intPtr(0)}); err != nil {
		t.Fatalf("Handle(first start_game) failed: %v", err)
	}
	if _, err := driver.Handle(&inbound.EndGame{Type: "end_game"}); err != nil {
		t.Fatalf("Handle(end_game) failed: %v", err)
	}
	if _, err := driver.Handle(&inbound.StartGame{Type: "start_game", ID: intPtr(1)}); err != nil {
		t.Fatalf("Handle(second start_game) failed: %v", err)
	}

	if agent.resets != 2 {
		t.Errorf("Reset calls = %d, want 2", agent.resets)
	}
}

func TestDriver_HandleEndGameMarksEnded(t *testing.T) {
	driver := mjairuntime.NewDriver("tsumogiri", "default", 0, ai.NewTsumogiriAgent(), nil)

	msg, err := driver.Handle(&inbound.EndGame{Type: "end_game"})
	if err != nil {
		t.Fatalf("Handle() failed: %v", err)
	}
	if msg != nil {
		t.Errorf("Handle() = %T, want nil", msg)
	}
	if !driver.Ended() {
		t.Error("Ended() = false, want true")
	}
}

func TestDriver_Ended(t *testing.T) {
	driver := mjairuntime.NewDriver("tsumogiri", "default", 0, ai.NewTsumogiriAgent(), nil)
	if driver.Ended() {
		t.Error("Ended() = true before start_game, want false")
	}

	if _, err := driver.Handle(&inbound.StartGame{Type: "start_game", ID: intPtr(0)}); err != nil {
		t.Fatalf("Handle(start_game) failed: %v", err)
	}
	if driver.Ended() {
		t.Error("Ended() = true after start_game, want false")
	}

	if _, err := driver.Handle(&inbound.EndGame{Type: "end_game"}); err != nil {
		t.Fatalf("Handle(end_game) failed: %v", err)
	}
	if !driver.Ended() {
		t.Error("Ended() = false after end_game, want true")
	}

	if _, err := driver.Handle(&inbound.StartGame{Type: "start_game", ID: intPtr(0)}); err != nil {
		t.Fatalf("Handle(second start_game) failed: %v", err)
	}
	if driver.Ended() {
		t.Error("Ended() = true after second start_game, want false")
	}
}

func TestDriver_HandleEventAfterEndGame(t *testing.T) {
	driver := mjairuntime.NewDriver("tsumogiri", "default", 0, ai.NewTsumogiriAgent(), nil)
	if _, err := driver.Handle(&inbound.StartGame{Type: "start_game", ID: intPtr(0)}); err != nil {
		t.Fatalf("Handle(start_game) failed: %v", err)
	}
	if _, err := driver.Handle(&inbound.EndGame{Type: "end_game"}); err != nil {
		t.Fatalf("Handle(end_game) failed: %v", err)
	}

	if _, err := driver.Handle(&inbound.Tsumo{Type: "tsumo", Actor: 0, Pai: "6m"}); err == nil {
		t.Error("Handle(tsumo) succeeded unexpectedly")
	}
}

func TestDriver_HandleEventBeforeStartGame(t *testing.T) {
	driver := mjairuntime.NewDriver("tsumogiri", "default", 0, ai.NewTsumogiriAgent(), nil)

	if _, err := driver.Handle(&inbound.Tsumo{Type: "tsumo", Actor: 0, Pai: "6m"}); err == nil {
		t.Fatal("Handle() succeeded unexpectedly")
	}
}

func TestDriver_HandleStartGameUsesFallbackIDWhenIDIsMissing(t *testing.T) {
	driver := mjairuntime.NewDriver("tsumogiri", "default", 2, ai.NewTsumogiriAgent(), nil)

	if _, err := driver.Handle(&inbound.StartGame{Type: "start_game"}); err != nil {
		t.Fatalf("Handle(start_game without id) failed: %v", err)
	}
}

func TestDriver_HandleStartGamePrefersMessageID(t *testing.T) {
	driver := mjairuntime.NewDriver("tsumogiri", "default", 4, ai.NewTsumogiriAgent(), nil)

	if _, err := driver.Handle(&inbound.StartGame{Type: "start_game", ID: intPtr(1)}); err != nil {
		t.Fatalf("Handle(start_game with id) failed: %v", err)
	}
}

func TestDriver_HandleStartGameDoesNotCarryMessageIDToNextGame(t *testing.T) {
	driver := mjairuntime.NewDriver("tsumogiri", "default", 4, ai.NewTsumogiriAgent(), nil)

	if _, err := driver.Handle(&inbound.StartGame{Type: "start_game", ID: intPtr(1)}); err != nil {
		t.Fatalf("Handle(first start_game) failed: %v", err)
	}
	if _, err := driver.Handle(&inbound.EndGame{Type: "end_game"}); err != nil {
		t.Fatalf("Handle(end_game) failed: %v", err)
	}
	if _, err := driver.Handle(&inbound.StartGame{Type: "start_game"}); err == nil {
		t.Fatal("Handle(second start_game without id) succeeded unexpectedly")
	}
}

func TestDriver_HandleStartGameRejectsInvalidFallbackID(t *testing.T) {
	driver := mjairuntime.NewDriver("tsumogiri", "default", 4, ai.NewTsumogiriAgent(), nil)

	if _, err := driver.Handle(&inbound.StartGame{Type: "start_game"}); err == nil {
		t.Fatal("Handle(start_game without id) succeeded unexpectedly")
	}
}

type recordingAgent struct {
	resets int
}

func (a *recordingAgent) Reset() {
	a.resets++
}

func (*recordingAgent) Decide(ai.Request) (ai.Decision, error) {
	return ai.Decision{}, nil
}

func intPtr(v int) *int {
	return &v
}
