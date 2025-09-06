package main

import (
	"fmt"
	"log"
	"os"
	"slices"

	"github.com/Apricot-S/mjai-manue-go/internal/ai"
	"github.com/Apricot-S/mjai-manue-go/internal/game"
	"github.com/Apricot-S/mjai-manue-go/internal/game/event/inbound"
	"github.com/Apricot-S/mjai-manue-go/internal/game/event/outbound"
	"github.com/Apricot-S/mjai-manue-go/internal/protocol/mjai"
	"github.com/Apricot-S/mjai-manue-go/tools/shared"
)

const playerName = "Manue014"

type Verifier struct {
	PlayerID int
	AI       ai.AI
	Decision outbound.Event
}

func NewVerifier(ai ai.AI) *Verifier {
	return &Verifier{AI: ai}
}

// ModifyAction is called before applying the action to the state.
// It can modify the action if needed.
func (v *Verifier) ModifyAction(action inbound.Event, g game.StateAnalyzer) (inbound.Event, error) {
	switch a := action.(type) {
	case *inbound.Error:
		return action, fmt.Errorf("error in the log: %+v", a)
	case *inbound.StartGame:
		v.PlayerID = slices.Index(a.Names[:], playerName)
		if v.PlayerID == -1 {
			return action, fmt.Errorf("player name %q not found in players %v", playerName, a.Names)
		}

		startGame, err := inbound.NewStartGame(v.PlayerID, a.Names)
		if err != nil {
			return action, err
		}
		action = startGame
	}

	decision, err := v.AI.DecideAction(g, v.PlayerID)
	if err != nil {
		return action, fmt.Errorf("failed to decide action: %w", err)
	}
	v.Decision = decision

	return action, nil
}

func (v *Verifier) VerifyAction(action inbound.Event, g game.StateViewer) (string, error) {
	switch a := action.(type) {
	case *inbound.Dahai:
		if a.Actor != v.PlayerID {
			return "", nil
		}
		da, ok := v.Decision.(*outbound.Dahai)
		if !ok {
			return "", fmt.Errorf("expected dahai decision, got %#v", v.Decision)
		}
		if a.Pai != da.Pai || a.Tsumogiri != da.Tsumogiri {
			return fmt.Sprintf("dahai mismatch:\nport:\n%+v\n\noriginal:\n%+v\n\n", da.Log, a), nil
		}
	case *inbound.Chi:
		if a.Actor != v.PlayerID {
			return "", nil
		}
		ch, ok := v.Decision.(*outbound.Chi)
		if !ok {
			return "", fmt.Errorf("expected chi decision, got %#v", v.Decision)
		}
		if a.Taken != ch.Taken || a.Consumed != ch.Consumed {
			return fmt.Sprintf("chi mismatch:\nport:\n%+v\n\noriginal:\n%+v\n\n", ch.Log, a), nil
		}
	case *inbound.Pon:
		if a.Actor != v.PlayerID {
			return "", nil
		}
		po, ok := v.Decision.(*outbound.Pon)
		if !ok {
			return "", fmt.Errorf("expected pon decision, got %#v", v.Decision)
		}
		if a.Target != po.Target || a.Taken != po.Taken || a.Consumed != po.Consumed {
			return fmt.Sprintf("pon mismatch:\nport:\n%+v\n\noriginal:\n%+v\n\n", po.Log, a), nil
		}
	case *inbound.Daiminkan:
		if a.Actor != v.PlayerID {
			return "", nil
		}
		dm, ok := v.Decision.(*outbound.Daiminkan)
		if !ok {
			return "", fmt.Errorf("expected daiminkan decision, got %#v", v.Decision)
		}
		if a.Target != dm.Target || a.Taken != dm.Taken {
			return fmt.Sprintf("daiminkan mismatch:\nport:\n%+v\n\noriginal:\n%+v\n\n", dm.Log, a), nil
		}
	case *inbound.Reach:
		if a.Actor != v.PlayerID {
			return "", nil
		}
		_, ok := v.Decision.(*outbound.Reach)
		if !ok {
			return "", fmt.Errorf("expected reach decision, got %#v", v.Decision)
		}
	case *inbound.Hora:
		if a.Actor != v.PlayerID {
			return "", nil
		}
		ho, ok := v.Decision.(*outbound.Hora)
		if !ok {
			return "", fmt.Errorf("expected hora decision, got %#v", v.Decision)
		}
		if a.Target != ho.Target || *a.Pai != ho.Pai {
			return fmt.Sprintf("hora mismatch:\nport:\n%+v\n\noriginal:\n%+v\n\n", ho.Log, a), nil
		}
	}
	return "", nil
}

func run(args []string) (bool, error) {
	paths, err := shared.GlobAll(args)
	if err != nil {
		return false, fmt.Errorf("error in glob: %w", err)
	}

	ai, err := ai.NewManueAIDefault()
	if err != nil {
		log.Fatalf("failed to create AI: %v", err)
	}
	verifier := NewVerifier(ai)

	archive := shared.NewArchive(paths, &mjai.MjaiAdapter{})
	hasMismatch := false
	onAction := func(action inbound.Event) error {
		action, err := verifier.ModifyAction(action, archive.StateAnalyzer())
		if err != nil {
			return err
		}
		if err := archive.StateUpdater().Update(action); err != nil {
			return err
		}
		detail, err := verifier.VerifyAction(action, archive.StateViewer())
		if err != nil {
			return err
		}
		if detail != "" {
			fmt.Println("VerifyAction mismatch:")
			fmt.Printf("state (after action):\n%s\n", archive.StateViewer().RenderBoard())
			fmt.Printf("%v\n", detail)
			hasMismatch = true
		}
		return nil
	}

	if err := archive.PlayLight(onAction); err != nil {
		return false, fmt.Errorf("error in processing log: %w", err)
	}

	return !hasMismatch, nil
}

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s <LOG_GLOB_PATTERNS>...\n", os.Args[0])
		os.Exit(2)
	}

	if ok, err := run(os.Args[1:]); err != nil {
		fmt.Println(err)
	} else if ok {
		fmt.Println("PASS")
	}
}
