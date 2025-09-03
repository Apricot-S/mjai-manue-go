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

// BeforeAction is called before applying the action to the state.
// It can modify the action if needed.
func (v *Verifier) BeforeAction(action inbound.Event, g game.StateAnalyzer) error {
	switch a := action.(type) {
	case *inbound.Error:
		return fmt.Errorf("error in the log: %+v", a)
	case *inbound.StartGame:
		v.PlayerID = slices.Index(a.Names[:], playerName)
		if v.PlayerID == -1 {
			return fmt.Errorf("player name %q not found in players %v", playerName, a.Names)
		}

		startGame, err := inbound.NewStartGame(v.PlayerID, a.Names)
		if err != nil {
			return err
		}
		action = startGame
	}

	decision, err := v.AI.DecideAction(g, v.PlayerID)
	if err != nil {
		return fmt.Errorf("failed to decide action: %w", err)
	}
	v.Decision = decision

	return nil
}

func (v *Verifier) VerifyAction(action inbound.Event, g game.StateViewer) error {
	return nil
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
	onAction := func(action inbound.Event) error {
		if err := verifier.BeforeAction(action, archive.StateAnalyzer()); err != nil {
			return err
		}
		if err := archive.StateUpdater().Update(action); err != nil {
			return err
		}
		if err := verifier.VerifyAction(action, archive.StateViewer()); err != nil {
			return err
		}
		return nil
	}

	if err := archive.PlayLight(onAction); err != nil {
		return false, fmt.Errorf("error in processing log: %w", err)
	}

	return true, nil
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
